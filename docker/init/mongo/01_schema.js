// 切换到借阅数据库
db = db.getSiblingDB('library_borrow');

// =============================================
// 创建集合
// =============================================

// 借阅记录集合
db.createCollection('borrow_records', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['book_id', 'user_id', 'borrow_date', 'due_date', 'status'],
            properties: {
                book_id: { bsonType: 'int' },
                book_title: { bsonType: 'string' },
                book_isbn: { bsonType: 'string' },
                user_id: { bsonType: 'int' },
                user_name: { bsonType: 'string' },
                borrow_date: { bsonType: 'date' },
                due_date: { bsonType: 'date' },
                return_date: { bsonType: 'date' },
                status: { enum: ['borrowed', 'returned', 'overdue', 'lost'] },
                renew_count: { bsonType: 'int' },
                fine: { bsonType: 'double' },
                fine_paid: { bsonType: 'bool' },
                operator_id: { bsonType: 'int' },
                remark: { bsonType: 'string' },
                created_at: { bsonType: 'date' },
                updated_at: { bsonType: 'date' }
            }
        }
    }
});

// 借阅规则集合
db.createCollection('borrow_rules', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['name', 'user_type'],
            properties: {
                name: { bsonType: 'string' },
                user_type: { bsonType: 'string' },
                max_books: { bsonType: 'int' },
                max_days: { bsonType: 'int' },
                max_renew_times: { bsonType: 'int' },
                fine_per_day: { bsonType: 'double' },
                max_fine: { bsonType: 'double' },
                description: { bsonType: 'string' },
                created_at: { bsonType: 'date' },
                updated_at: { bsonType: 'date' }
            }
        }
    }
});

// 预约记录集合
db.createCollection('reservations', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['book_id', 'user_id', 'status'],
            properties: {
                book_id: { bsonType: 'int' },
                book_title: { bsonType: 'string' },
                user_id: { bsonType: 'int' },
                user_name: { bsonType: 'string' },
                status: { enum: ['waiting', 'notified', 'completed', 'cancelled'] },
                queue_position: { bsonType: 'int' },
                notify_date: { bsonType: 'date' },
                expire_date: { bsonType: 'date' },
                created_at: { bsonType: 'date' },
                updated_at: { bsonType: 'date' }
            }
        }
    }
});

// 借阅统计集合
db.createCollection('borrow_statistics', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['date'],
            properties: {
                date: { bsonType: 'date' },
                total_borrowed: { bsonType: 'int' },
                total_returned: { bsonType: 'int' },
                total_overdue: { bsonType: 'int' },
                total_fine: { bsonType: 'double' },
                created_at: { bsonType: 'date' }
            }
        }
    }
});

// =============================================
// 创建索引
// =============================================

// 借阅记录索引
db.borrow_records.createIndex({ user_id: 1, status: 1 });
db.borrow_records.createIndex({ book_id: 1, status: 1 });
db.borrow_records.createIndex({ status: 1, due_date: 1 });
db.borrow_records.createIndex({ borrow_date: -1 });
db.borrow_records.createIndex({ created_at: -1 });

// 预约记录索引
db.reservations.createIndex({ book_id: 1, status: 1 });
db.reservations.createIndex({ user_id: 1, status: 1 });
db.reservations.createIndex({ status: 1, created_at: 1 });

// 借阅统计索引
db.borrow_statistics.createIndex({ date: 1 }, { unique: true });

// =============================================
// 插入默认借阅规则
// =============================================

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
    },
    {
        name: '普通用户借阅规则',
        user_type: 'user',
        max_books: 3,
        max_days: 30,
        max_renew_times: 2,
        fine_per_day: 0.5,
        max_fine: 30,
        description: '普通用户最多可借3本书，借期30天，可续借2次',
        created_at: new Date(),
        updated_at: new Date()
    }
]);

// =============================================
// 插入测试借阅记录
// =============================================

var now = new Date();
var borrowDate = new Date(now.getTime() - 15 * 24 * 60 * 60 * 1000); // 15天前
var dueDate = new Date(now.getTime() + 15 * 24 * 60 * 60 * 1000); // 15天后

db.borrow_records.insertMany([
    {
        book_id: 1,
        book_title: '深入理解计算机系统',
        book_isbn: '978-7-111-54784-2',
        user_id: 2,
        user_name: '张三',
        borrow_date: borrowDate,
        due_date: dueDate,
        status: 'borrowed',
        renew_count: 0,
        fine: 0,
        fine_paid: false,
        created_at: borrowDate,
        updated_at: borrowDate
    },
    {
        book_id: 2,
        book_title: 'JavaScript高级程序设计',
        book_isbn: '978-7-115-42533-4',
        user_id: 3,
        user_name: '李四',
        borrow_date: new Date(now.getTime() - 35 * 24 * 60 * 60 * 1000),
        due_date: new Date(now.getTime() - 5 * 24 * 60 * 60 * 1000),
        status: 'overdue',
        renew_count: 1,
        fine: 2.5,
        fine_paid: false,
        created_at: new Date(now.getTime() - 35 * 24 * 60 * 60 * 1000),
        updated_at: now
    },
    {
        book_id: 3,
        book_title: 'Python编程从入门到实践',
        book_isbn: '978-7-115-42602-7',
        user_id: 4,
        user_name: '王五',
        borrow_date: new Date(now.getTime() - 20 * 24 * 60 * 60 * 1000),
        due_date: new Date(now.getTime() - 5 * 24 * 60 * 60 * 1000),
        return_date: new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000),
        status: 'returned',
        renew_count: 0,
        fine: 1.0,
        fine_paid: true,
        created_at: new Date(now.getTime() - 20 * 24 * 60 * 60 * 1000),
        updated_at: new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000)
    }
]);

// =============================================
// 创建视图
// =============================================

// 当前借阅中的记录视图
db.createView('current_borrows', 'borrow_records', [
    { $match: { status: { $in: ['borrowed', 'overdue'] } } },
    { $sort: { borrow_date: -1 } }
]);

// 逾期记录视图
db.createView('overdue_borrows', 'borrow_records', [
    { $match: { status: 'overdue' } },
    { $sort: { due_date: 1 } }
]);

// =============================================
// 创建函数（存储过程）
// =============================================

// 计算罚款函数
db.system.js.save({
    _id: 'calculateFine',
    value: function(borrowId) {
        var record = db.borrow_records.findOne({ _id: borrowId });
        if (!record || record.status === 'returned') return 0;
        
        var now = new Date();
        if (now <= record.due_date) return 0;
        
        var days = Math.floor((now - record.due_date) / (24 * 60 * 60 * 1000));
        var rule = db.borrow_rules.findOne({ user_type: 'user' });
        var finePerDay = rule ? rule.fine_per_day : 0.5;
        
        return days * finePerDay;
    }
});

// 检查并更新逾期状态函数
db.system.js.save({
    _id: 'updateOverdueStatus',
    value: function() {
        var now = new Date();
        var result = db.borrow_records.updateMany(
            { 
                status: 'borrowed', 
                due_date: { $lt: now } 
            },
            { 
                $set: { 
                    status: 'overdue',
                    updated_at: now
                } 
            }
        );
        return result.modifiedCount;
    }
});

// 完成
print('MongoDB initialization completed!');
print('Collections: ' + db.getCollectionNames().join(', '));
