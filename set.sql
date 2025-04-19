-- users 테이블
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- work_logs 테이블
CREATE TABLE IF NOT EXISTS work_logs (
  id         SERIAL    PRIMARY KEY,
  user_id    INTEGER   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content    TEXT      NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- work_logs의 하루 하나 보장하려면 함수식 인덱스 추가
CREATE UNIQUE INDEX IF NOT EXISTS
  idx_work_logs_user_date
ON work_logs (
  user_id,
  (created_at::date)
);
