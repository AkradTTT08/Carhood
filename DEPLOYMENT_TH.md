# คู่มือการติดตั้งระบบ Carhood (TTT Space Quiz) บน Linux Server

คู่มือนี้จะอธิบายขั้นตอนการนำระบบ Carhood ไปรันบน Server จริงโดยใช้ **Docker** และ **Docker Compose** แบบละเอียดทุกขั้นตอน

---

## 📋 สิ่งที่ต้องเตรียมก่อนเริ่ม (Prerequisites)

1. **Linux Server**: แนะนำเป็น Ubuntu 22.04 LTS หรือเวอร์ชันที่ใกล้เคียง
2. **Docker**: ติดตั้ง Docker Engine เรียบร้อยแล้ว
3. **Docker Compose**: ติดตั้ง Docker Compose เรียบร้อยแล้ว (V2 ขึ้นไป)
4. **Git**: สำหรับใช้ในการ Clone Code จาก Repository

---

## 🛡️ การตั้งค่า Firewall และ Port ที่ต้องเปิด

เพื่อให้ผู้ใช้งานภายนอกสามารถเข้าถึงระบบได้ คุณต้องเปิด Port ต่อไปนี้ที่ Firewall (เช่น ufw หรือ Cloud Security Group):

- **Port 3000 (TCP)**: สำหรับหน้าเว็บหลัก (Frontend)
- **Port 8081 (TCP)**: สำหรับ API และ WebSocket (Backend)
- **Port 22 (TCP)**: สำหรับ SSH (ปกติเปิดอยู่แล้วสำหรับจัดการ Server)

*หมายเหตุ: สำหรับ Port 27017 (MongoDB) แนะนำให้ปิดไว้สำหรับใช้งานภายในเท่านั้นเพื่อความปลอดภัย*

---

## 🚀 ขั้นตอนการติดตั้ง (Deployment Steps)

### 1. การนำ Code ขึ้น Server

วิธีที่แนะนำที่สุดคือการใช้งาน **Git** โดยคุณสามารถเลือกทำได้ 2 วิธี:

#### วิธีที่ A: ผ่าน HTTPS (ง่ายที่สุด)
```bash
git clone https://github.com/your-username/Carhood.git
cd Carhood
```

#### วิธีที่ B: ผ่าน SSH
หากเป็น Private Repo แนะนำให้ตั้งค่า SSH Key ก่อน:
1. สร้าง SSH Key: `ssh-keygen -t ed25519 -C "your_email@example.com"`
2. Copy Key ไปใส่ใน GitHub/GitLab: `cat ~/.ssh/id_ed25519.pub`
3. Clone: `git clone git@github.com:your-username/Carhood.git`

---

### 2. สั่งรันระบบด้วย Docker Compose
ใช้คำสั่งเดี่ยวเพื่อ Build และ Start ทุก Service:

```bash
docker-compose up -d --build
```

---

### 3. การตรวจสอบสถานะ
```bash
docker-compose ps
docker-compose logs -f
```

---

## 🌐 การเข้าใช้งานระบบ

เมื่อระบบรันเสร็จแล้ว สามารถเข้าใช้งานได้ผ่าน URL ต่อไปนี้:

- **Frontend (หน้าเว็บ)**: `http://<IP-SERVER>:3000`
- **Backend API**: `http://<IP-SERVER>:8081`

---

## 🛠️ การแก้ไขปัญหาเบื้องต้น (Troubleshooting)

| ปัญหา | วิธีแก้ไข |
| :--- | :--- |
| **เข้าหน้าเว็บไม่ได้** | ตรวจสอบว่าเปิด Port 3000 ที่ Firewall หรือยัง |
| **เชื่อมต่อ WebSocket ไม่ได้** | ตรวจสอบว่าเปิด Port 8081 ที่ Firewall หรือยัง |
| **ต้องการหยุดการทำงาน** | ใช้คำสั่ง `docker-compose down` |
| **ต้องการอัปเดตโค้ดใหม่** | `git pull origin main` แล้วรัน `docker-compose up -d --build` อีกครั้ง |
