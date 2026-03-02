CREATE TABLE email_notification(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    subject VARCHAR(255) NOT NULL,
    text TEXT NOT NULL,
    send_at DATETIME NOT NULL,
    status ENUM('pending', 'queued','sent', 'failed') DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_status_send_at (status, send_at)
);