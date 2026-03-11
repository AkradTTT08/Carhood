# คู่มือการติดตั้งระบบ Carhood (TTT Space Quiz) บน Linux Server

คู่มือนี้จะอธิบายขั้นตอนการนำระบบ Carhood ไปรันบน Server จริงโดยใช้ **Docker** และ **Docker Compose** แบบละเอียดทุกขั้นตอน

---

## 📋 สิ่งที่ต้องเตรียมก่อนเริ่ม (Prerequisites)

1. **Linux Server**: แนะนำเป็น Ubuntu 22.04 LTS หรือเวอร์ชันที่ใกล้เคียง
2. **Docker**: ติดตั้ง Docker Engine เรียบร้อยแล้ว
3. **Docker Compose**: ติดตั้ง Docker Compose เรียบร้อยแล้ว (V2 ขึ้นไป)
4. **Git**: สำหรับใช้ในการ Clone Code จาก Repository

---

## 🚀 ขั้นตอนการติดตั้ง (Deployment Steps)

### 1. การนำ Code ขึ้น Server

วิธีที่แนะนำที่สุดคือการใช้งาน **Git** โดยคุณสามารถเลือกทำได้ 2 วิธี:

#### วิธีที่ A: ผ่าน HTTPS (ง่ายที่สุด)
```bash
git clone https://github.com/your-username/Carhood.git
cd Carhood
```

#### วิธีที่ B: ผ่าน SSH (แนะนำสำหรับการจัดการระยะยาว)
หากเป็น Private Repo แนะนำให้ตั้งค่า SSH Key ก่อน:
1. สร้าง SSH Key: `ssh-keygen -t ed25519 -C "your_email@example.com"`
2. Copy Key ไปใส่ใน GitHub/GitLab: `cat ~/.ssh/id_ed25519.pub`
3. Clone: `git clone git@github.com:your-username/Carhood.git`

---

### 2. โครงสร้างโปรเจกต์ (Project Structure)
ตรวจสอบว่าคุณอยู่ในโฟลเดอร์หลักของโปรเจกต์ ซึ่งควรมีโครงสร้างดังนี้:
- `backend/` - โค้ดส่วนหลังบ้าน (Go)
- `frontend/` - โค้ดส่วนหน้าบ้าน (SvelteKit)
- `docker-compose.yml` - ไฟล์ตั้งค่าสำหรับรัน Docker
- `README.md` และ `DEPLOYMENT.md`

---

### 3. การตั้งค่า Environment (Optional)
โดยปกติระบบถูกตั้งค่าให้รันได้ทันที แต่หากต้องการเปลี่ยน Database URL สามารถสร้างไฟล์ `.env` ไว้ที่ Root:
```env
MONGODB_URI=mongodb://mongodb:27017
```

---

### 4. สั่งรันระบบด้วย Docker Compose
ใช้คำสั่งเดียวเพื่อ Build และ Start ทุก Service (Backend, Frontend, MongoDB):

```bash
docker-compose up -d --build
```

**คำอธิบายคำสั่ง:**
- `up`: สั่งให้ Docker เริ่มทำงานตามไฟล์ compose
- `-d`: รันใน Background mode (ปล่อยให้ระบบทำงานต่อแม้เราจะปิด Terminal)
- `--build`: สั่งให้ Build Image ใหม่เสมอ (ป้องกันการใช้ Cache เก่าที่โค้ดยังไม่อัปเดต)

---

### 5. การตรวจสอบสถานะ
ตรวจสอบว่าคอนเทนเนอร์ทำงานปกติหรือไม่:
```bash
docker-compose ps
```

ดู Log การทำงาน (เพื่อเช็ค Error):
```bash
docker-compose logs -f
```

---

## 🌐 การเข้าใช้งานระบบ

เมื่อระบบรันเสร็จแล้ว สามารถเข้าใช้งานได้ผ่าน URL ต่อไปนี้:

- **Frontend (หน้าเว็บ)**: `http://<IP-SERVER>:3000`
- **Backend API**: `http://<IP-SERVER>:8081`
- **MongoDB**: จะทำงานภายในเฉพาะเครื่อง (Internal Network)

---

## 🛠️ การแก้ไขปัญหาเบื้องต้น (Troubleshooting)

| ปัญหา | วิธีแก้ไข |
| :--- | :--- |
| **รันไม่ได้ หรือ Error Port 3000** | ตรวจสอบว่ามี Service อื่นรันค้างอยู่บน Port 3000 หรือไม่ |
| **Backend เชื่อมต่อ DB ไม่ได้** | รอสักครู่เพื่อให้ MongoDB เริ่มต้นทำงานเสร็จ หรือเช็ค `docker-compose logs backend` |
| **ต้องการหยุดการทำงาน** | ใช้คำสั่ง `docker-compose down` |
| **ต้องการอัปเดตโค้ดใหม่** | `git pull origin main` แล้วรัน `docker-compose up -d --build` อีกครั้ง |

---

## 💡 หมายเหตุเพิ่มเติม
- **โฟลเดอร์อัปโหลด**: รูปภาพที่อัปโหลดจะถูกเก็บไว้ในโฟลเดอร์ `uploads` ใน Backend Container ซึ่งถูกตั้งค่าให้เป็น Volume ที่คงทน (Persistent)
- **URL อัตโนมัติ**: Frontend ถูกออกแบบมาให้ตรวจหา IP ของ Server เองโดยอัตโนมัติ ไม่จำเป็นต้องแก้โค้ดเมื่อเปลี่ยน Server
