import { ec as EC } from 'elliptic';
import crypto from 'crypto';

const ec = new EC('secp256k1');

interface EncryptedPackage {
  ephemeralPubKey: string;
  iv: string;
  authTag: string;
  ciphertext: string;
}

/**
 * Encrypts a plaintext string (Alice's seed) using Bob's Secp256k1 Public Key.
 * This utilizes a custom ECIES structure optimized for 48-hour hackathon execution velocity.
 */
export function encryptPayload(plaintext: string, beneficiaryPubKeyHex: string): string {
  try {
    // 1. Generate an ephemeral, single-use public/private keypair
    const ephemeralKey = ec.genKeyPair();
    
    // 2. Parse Bob's public key from hex format
    const targetPubKey = ec.keyFromPublic(beneficiaryPubKeyHex, 'hex');
    
    // 3. Perform a Diffie-Hellman Key Exchange (ECDH) to derive a shared secret key
    const sharedSecretPoint = ephemeralKey.derive(targetPubKey.getPublic());
    const sharedSecretBuffer = Buffer.from(sharedSecretPoint.toString(16, 64), 'hex');
    
    // Hash the shared secret with SHA-256 to ensure a uniform 256-bit encryption key
    const aesKey = crypto.createHash('sha256').update(sharedSecretBuffer).digest();
    
    // 4. Encrypt the plaintext payload using enterprise-standard AES-256-GCM
    const iv = crypto.randomBytes(12); // Initialization vector
    const cipher = crypto.createCipheriv('aes-256-gcm', aesKey, iv);
    
    let ciphertext = cipher.update(plaintext, 'utf8', 'hex');
    ciphertext += cipher.final('hex');
    
    const authTag = cipher.getAuthTag().toString('hex');
    
    // 5. Package the resulting components together
    const finalPackage: EncryptedPackage = {
      ephemeralPubKey: ephemeralKey.getPublic('hex'),
      iv: iv.toString('hex'),
      authTag,
      ciphertext: ciphertext
    };
    
    // Convert the payload object to a standard string for database transit
    return btoa(JSON.stringify(finalPackage));
  } catch (error) {
    throw new Error(`Browser encryption failed: ${(error as Error).message}`);
  }
}

/**
 * Decrypts a ciphertext string package using Bob's private key.
 * Bob uses this utility locally inside his browser when pulling a dormant vault.
 */
export function decryptPayload(base64Package: string, beneficiaryPrivateKeyHex: string): string {
  try {
    // Decode the Base64 transit string back into individual crypto parameters
    const decodedPackage: EncryptedPackage = JSON.parse(atob(base64Package));
    
    const beneficiaryKey = ec.keyFromPrivate(beneficiaryPrivateKeyHex, 'hex');
    const ephemeralPubKey = ec.keyFromPublic(decodedPackage.ephemeralPubKey, 'hex');
    
    // Re-derive the exact same shared secret point using Bob's private key
    const sharedSecretPoint = beneficiaryKey.derive(ephemeralPubKey.getPublic());
    const sharedSecretBuffer = Buffer.from(sharedSecretPoint.toString(16, 64), 'hex');
    const aesKey = crypto.createHash('sha256').update(sharedSecretBuffer).digest();
    
    const iv = Buffer.from(decodedPackage.iv, 'hex');
    const authTag = Buffer.from(decodedPackage.authTag, 'hex');
    
    const decipher = crypto.createDecipheriv('aes-256-gcm', aesKey, iv);
    decipher.setAuthTag(authTag);
    
    let decrypted = decipher.update(decodedPackage.ciphertext, 'hex', 'utf8');
    decrypted += decipher.final('utf8');
    
    return decrypted;
  } catch (error) {
    throw new Error(`Local decryption failed: Authentication tag mismatch or invalid private key.`);
  }
}