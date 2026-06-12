# Continuum

A privacy-preserving Bitcoin inheritance protocol that uses **Lightning Network proof-of-life payments** to determine wallet activity and trigger automated asset recovery when a user becomes inactive.

---

## ⚡ Overview

**Continuum** is a decentralized inheritance system built on Bitcoin and Lightning infrastructure.

It solves a critical problem in crypto:

> Millions of Bitcoin are permanently lost due to forgotten keys and lack of inheritance mechanisms.

Finality introduces a **proof-of-life mechanism** using Lightning payments to ensure digital assets can be safely passed on when the owner is no longer active.

---

## 🧠 Core Idea

Instead of relying on legal systems, emails, or centralized identity verification:

- Users prove they are alive by sending a Lightning payment (“ping”)
- The system tracks activity over time
- If pings stop beyond a defined threshold, recovery is triggered
- Beneficiaries receive access to encrypted recovery instructions

---

## 🔐 Key Features

### ⚡ Lightning Proof-of-Life
Users periodically send small Lightning payments to confirm wallet activity.

### 🧾 Inheritance Vaults
Users create vaults that define:
- Beneficiaries (via public keys)
- Inactivity threshold
- Recovery rules

### 🔒 Encrypted Recovery Packages
Sensitive recovery data is encrypted and only accessible to beneficiaries.

### ⏱ Continuum Engine
Automated scheduler monitors inactivity and triggers inheritance flow.

### 👥 Multi-Beneficiary Support
Multiple recipients can be assigned roles (primary / backup).

---

## 🏗 System Architecture
