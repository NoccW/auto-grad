-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    openId VARCHAR(64) NOT NULL UNIQUE,
    name TEXT,
    email VARCHAR(320),
    loginMethod VARCHAR(64),
    role ENUM('user', 'admin') DEFAULT 'user' NOT NULL,
    userRole ENUM('parent', 'teacher'),
    passwordHash TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    lastSignedIn TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    INDEX idx_openId (openId),
    INDEX idx_userRole (userRole),
    INDEX idx_createdAt (createdAt)
);

-- 创建批改记录表
CREATE TABLE IF NOT EXISTS grading_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId INT NOT NULL,
    paperImageUrl TEXT NOT NULL,
    answerImageUrl TEXT,
    ocrResult TEXT,
    aiScore INT,
    wrongQuestions TEXT,
    correctAnswers TEXT,
    status ENUM('pending', 'processing', 'completed', 'failed') DEFAULT 'pending' NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    INDEX idx_userId (userId),
    INDEX idx_status (status),
    INDEX idx_createdAt (createdAt),
    FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建教师任务表（预留功能）
CREATE TABLE IF NOT EXISTS teacher_tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId INT NOT NULL,
    targetUrl VARCHAR(500) NOT NULL,
    account VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    taskId VARCHAR(255) NOT NULL,
    status ENUM('pending', 'running', 'completed', 'failed') DEFAULT 'pending' NOT NULL,
    progress VARCHAR(50),
    totalPapers INT,
    completedPapers INT DEFAULT 0,
    errorMessage TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    completedAt TIMESTAMP NULL,
    INDEX idx_userId (userId),
    INDEX idx_status (status),
    INDEX idx_createdAt (createdAt),
    FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);