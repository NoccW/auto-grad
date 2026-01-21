-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    openId VARCHAR(64) NOT NULL UNIQUE,
    name TEXT,
    email VARCHAR(320),
    loginMethod VARCHAR(64),
    role VARCHAR(191) DEFAULT 'user',
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
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT UNSIGNED NOT NULL,
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

-- 创建教师任务表
CREATE TABLE IF NOT EXISTS teacher_tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT UNSIGNED NOT NULL,
    targetUrl VARCHAR(500) NOT NULL,
    account VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    taskId VARCHAR(255) NOT NULL,
    status ENUM('pending', 'running', 'completed', 'failed', 'cancelled') DEFAULT 'pending' NOT NULL,
    progress VARCHAR(50),
    totalPapers INT DEFAULT 0,
    completedPapers INT DEFAULT 0,
    failedPapers INT DEFAULT 0,
    averageScore DECIMAL(5,2) DEFAULT 0,
    errorMessage TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    completedAt TIMESTAMP NULL,
    INDEX idx_userId (userId),
    INDEX idx_status (status),
    INDEX idx_createdAt (createdAt),
    FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建教师任务执行记录表
CREATE TABLE IF NOT EXISTS teacher_task_executions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    taskId BIGINT UNSIGNED NOT NULL,
    paperId VARCHAR(255),
    studentName VARCHAR(255),
    score INT DEFAULT 0,
    ocrResult TEXT,
    aiFeedback TEXT,
    status ENUM('pending', 'processing', 'completed', 'failed') DEFAULT 'pending',
    errorMessage TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    INDEX idx_taskId (taskId),
    INDEX idx_status (status),
    FOREIGN KEY (taskId) REFERENCES teacher_tasks(id) ON DELETE CASCADE
);

-- 创建任务执行统计表
CREATE TABLE IF NOT EXISTS task_statistics (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    taskId BIGINT UNSIGNED NOT NULL UNIQUE,
    totalPapers INT DEFAULT 0,
    completedPapers INT DEFAULT 0,
    failedPapers INT DEFAULT 0,
    averageScore DECIMAL(5,2) DEFAULT 0,
    maxScore INT DEFAULT 0,
    minScore INT DEFAULT 100,
    passRate DECIMAL(5,2) DEFAULT 0,
    excellenceRate DECIMAL(5,2) DEFAULT 0,
    scoreDistribution JSON,
    errorDistribution JSON,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    INDEX idx_taskId (taskId),
    FOREIGN KEY (taskId) REFERENCES teacher_tasks(id) ON DELETE CASCADE
);