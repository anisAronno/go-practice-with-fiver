#!/bin/bash
docker exec gofiver-api sh -c "cat <<'EOF' | mysql -h mysql -u root -pbs@123 gofiver
CREATE INDEX IF NOT EXISTS idx_blogs_id_desc ON blogs(id DESC);
CREATE INDEX IF NOT EXISTS idx_blogs_user_id ON blogs(user_id);
CREATE INDEX IF NOT EXISTS idx_blogs_deleted ON blogs(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_deleted ON users(deleted_at);
ANALYZE TABLE blogs;
ANALYZE TABLE users;
EOF"
echo "Indexes created successfully"
