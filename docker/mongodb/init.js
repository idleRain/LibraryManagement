// 创建借阅数据库
db = db.getSiblingDB('library_borrow');

// 创建借阅记录集合
db.createCollection('borrow_records');

// 创建借阅规则集合
db.createCollection('borrow_rules');

// 创建预约记录集合
db.createCollection('reservations');

// 创建索引
db.borrow_records.createIndex({ user_id: 1, status: 1 });
db.borrow_records.createIndex({ book_id: 1, status: 1 });
db.borrow_records.createIndex({ due_date: 1 });
db.borrow_records.createIndex({ created_at: -1 });

db.borrow_rules.createIndex({ user_type: 1 });

db.reservations.createIndex({ user_id: 1, status: 1 });
db.reservations.createIndex({ book_id: 1, status: 1 });

// 插入默认借阅规则
db.borrow_rules.insertMany([
  {
    name: '学生借阅规则',
    user_type: 'student',
    max_books: 5,
    max_days: 30,
    max_renew_times: 2,
    fine_per_day: 0.5,
    max_fine: 50,
    description: '学生最多可借5本书，借期30天，可续借2次',
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: '教师借阅规则',
    user_type: 'teacher',
    max_books: 10,
    max_days: 60,
    max_renew_times: 3,
    fine_per_day: 0.3,
    max_fine: 30,
    description: '教师最多可借10本书，借期60天，可续借3次',
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: '职工借阅规则',
    user_type: 'staff',
    max_books: 3,
    max_days: 14,
    max_renew_times: 1,
    fine_per_day: 0.5,
    max_fine: 20,
    description: '职工最多可借3本书，借期14天，可续借1次',
    created_at: new Date(),
    updated_at: new Date()
  }
]);

print('MongoDB initialization completed!');
