Go Mysql Gateway

merupakan Api gateway untuk mengakses database mysql.
setiap request akan dicatat dalam log pada table logs

Stack :
1. Go
2. Gorm
3. Viper
4. Fiber

setelah running pastikan table logs telah terbentuk dalam database yang telah didaftarkan
dan sesuaikan dengan struktur berikut :

CREATE TABLE LOGS (
    id INT AUTO_INCREMENT,
    LEVEL VARCHAR(10),
    message TEXT,
    created_at DATETIME,
    PRIMARY KEY (id, created_at)
) ENGINE=INNODB PARTITION BY RANGE (TO_DAYS(created_at)) (
    PARTITION p0 VALUES LESS THAN (TO_DAYS('2023-01-01')),
    PARTITION p1 VALUES LESS THAN (TO_DAYS('2023-02-01')),
    PARTITION p2 VALUES LESS THAN (TO_DAYS('2023-03-01')),
    PARTITION p3 VALUES LESS THAN (TO_DAYS('2023-04-01')),
    PARTITION p4 VALUES LESS THAN (TO_DAYS('2023-05-01')),
    PARTITION p5 VALUES LESS THAN (TO_DAYS('2023-06-01')),
    PARTITION p6 VALUES LESS THAN (TO_DAYS('2023-07-01')),
    PARTITION p7 VALUES LESS THAN (TO_DAYS('2023-08-01')),
    PARTITION p8 VALUES LESS THAN (TO_DAYS('2023-09-01')),
    PARTITION p9 VALUES LESS THAN (TO_DAYS('2023-10-01')),
    PARTITION p10 VALUES LESS THAN (TO_DAYS('2023-11-01')),
    PARTITION p11 VALUES LESS THAN (TO_DAYS('2023-12-01'))
);
CREATE INDEX idx_created_at ON LOGS (created_at);

go build -ldflags="-s -w" main.go
