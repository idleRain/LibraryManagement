-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS library_system 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE library_system;

-- =============================================
-- 用户表
-- =============================================
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    real_name VARCHAR(50) COMMENT '真实姓名',
    avatar VARCHAR(255) COMMENT '头像URL',
    status TINYINT DEFAULT 1 COMMENT '状态: 1启用 0禁用',
    last_login_at TIMESTAMP NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) COMMENT '最后登录IP',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- =============================================
-- 角色表
-- =============================================
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL COMMENT '角色名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '角色编码',
    description VARCHAR(255) COMMENT '角色描述',
    sort_order INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态: 1启用 0禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_code (code),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- =============================================
-- 权限表
-- =============================================
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    resource VARCHAR(50) NOT NULL COMMENT '资源名称',
    action VARCHAR(20) NOT NULL COMMENT '操作: create, read, update, delete, *',
    parent_id BIGINT UNSIGNED DEFAULT 0 COMMENT '父权限ID',
    description VARCHAR(255) COMMENT '权限描述',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_code (code),
    INDEX idx_resource (resource),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- =============================================
-- 用户角色关联表
-- =============================================
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_role (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- =============================================
-- 角色权限关联表
-- =============================================
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    role_id BIGINT UNSIGNED NOT NULL,
    permission_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- =============================================
-- 图书表
-- =============================================
CREATE TABLE IF NOT EXISTS books (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    isbn VARCHAR(20) NOT NULL UNIQUE COMMENT 'ISBN编号',
    title VARCHAR(200) NOT NULL COMMENT '书名',
    subtitle VARCHAR(200) COMMENT '副标题',
    author VARCHAR(100) COMMENT '作者',
    translator VARCHAR(100) COMMENT '译者',
    publisher VARCHAR(100) COMMENT '出版社',
    publish_date DATE COMMENT '出版日期',
    edition VARCHAR(50) COMMENT '版次',
    category VARCHAR(50) COMMENT '分类',
    language VARCHAR(20) DEFAULT '中文' COMMENT '语言',
    pages INT COMMENT '页数',
    price DECIMAL(10,2) COMMENT '价格',
    cover_image VARCHAR(255) COMMENT '封面图片',
    description TEXT COMMENT '简介',
    keywords VARCHAR(255) COMMENT '关键词',
    status TINYINT DEFAULT 1 COMMENT '状态: 1上架 0下架',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_isbn (isbn),
    INDEX idx_title (title),
    INDEX idx_author (author),
    INDEX idx_category (category),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图书表';

-- =============================================
-- 库存表
-- =============================================
CREATE TABLE IF NOT EXISTS stocks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    book_id BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '图书ID',
    total_quantity INT DEFAULT 0 COMMENT '总库存',
    available_quantity INT DEFAULT 0 COMMENT '可借数量',
    borrowed_quantity INT DEFAULT 0 COMMENT '借出数量',
    sold_quantity INT DEFAULT 0 COMMENT '已售数量',
    damaged_quantity INT DEFAULT 0 COMMENT '损坏数量',
    location VARCHAR(100) COMMENT '存放位置',
    warning_threshold INT DEFAULT 5 COMMENT '库存预警阈值',
    last_stock_in_at TIMESTAMP NULL COMMENT '最后入库时间',
    last_stock_out_at TIMESTAMP NULL COMMENT '最后出库时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    INDEX idx_book_id (book_id),
    INDEX idx_available (available_quantity),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存表';

-- =============================================
-- 库存变动记录表
-- =============================================
CREATE TABLE IF NOT EXISTS stock_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    book_id BIGINT UNSIGNED NOT NULL COMMENT '图书ID',
    type ENUM('in', 'out', 'adjust') NOT NULL COMMENT '类型: 入库/出库/调整',
    quantity INT NOT NULL COMMENT '变动数量',
    before_quantity INT COMMENT '变动前数量',
    after_quantity INT COMMENT '变动后数量',
    reason VARCHAR(255) COMMENT '原因',
    related_order_type VARCHAR(50) COMMENT '关联订单类型',
    related_order_id BIGINT UNSIGNED COMMENT '关联订单ID',
    operator_id BIGINT UNSIGNED COMMENT '操作人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (operator_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_book_id (book_id),
    INDEX idx_type (type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存变动记录表';

-- =============================================
-- 供应商表
-- =============================================
CREATE TABLE IF NOT EXISTS suppliers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '供应商名称',
    code VARCHAR(50) UNIQUE COMMENT '供应商编码',
    contact VARCHAR(50) COMMENT '联系人',
    phone VARCHAR(20) COMMENT '电话',
    email VARCHAR(100) COMMENT '邮箱',
    address VARCHAR(255) COMMENT '地址',
    bank_name VARCHAR(100) COMMENT '开户银行',
    bank_account VARCHAR(50) COMMENT '银行账号',
    tax_number VARCHAR(50) COMMENT '税号',
    description TEXT COMMENT '备注',
    status TINYINT DEFAULT 1 COMMENT '状态: 1启用 0禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_name (name),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='供应商表';

-- =============================================
-- 采购订单表
-- =============================================
CREATE TABLE IF NOT EXISTS purchase_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE COMMENT '订单号',
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT '供应商ID',
    total_amount DECIMAL(12,2) DEFAULT 0 COMMENT '总金额',
    paid_amount DECIMAL(12,2) DEFAULT 0 COMMENT '已付金额',
    status ENUM('pending', 'approved', 'receiving', 'completed', 'cancelled') DEFAULT 'pending' COMMENT '状态',
    operator_id BIGINT UNSIGNED COMMENT '操作人ID',
    approved_by BIGINT UNSIGNED COMMENT '审核人ID',
    approved_at TIMESTAMP NULL COMMENT '审核时间',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE RESTRICT,
    FOREIGN KEY (operator_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_order_no (order_no),
    INDEX idx_supplier_id (supplier_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='采购订单表';

-- =============================================
-- 采购订单明细表
-- =============================================
CREATE TABLE IF NOT EXISTS purchase_order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    book_id BIGINT UNSIGNED NOT NULL COMMENT '图书ID',
    quantity INT NOT NULL COMMENT '采购数量',
    unit_price DECIMAL(10,2) NOT NULL COMMENT '单价',
    total_price DECIMAL(10,2) NOT NULL COMMENT '总价',
    received_quantity INT DEFAULT 0 COMMENT '已收货数量',
    status ENUM('pending', 'partial', 'completed') DEFAULT 'pending' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE RESTRICT,
    INDEX idx_order_id (order_id),
    INDEX idx_book_id (book_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='采购订单明细表';

-- =============================================
-- 销售订单表
-- =============================================
CREATE TABLE IF NOT EXISTS sale_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE COMMENT '订单号',
    customer_name VARCHAR(50) COMMENT '客户姓名',
    customer_phone VARCHAR(20) COMMENT '客户电话',
    customer_email VARCHAR(100) COMMENT '客户邮箱',
    total_amount DECIMAL(12,2) DEFAULT 0 COMMENT '订单金额',
    discount_amount DECIMAL(12,2) DEFAULT 0 COMMENT '优惠金额',
    pay_amount DECIMAL(12,2) DEFAULT 0 COMMENT '实付金额',
    payment_method ENUM('cash', 'wechat', 'alipay', 'card', 'other') COMMENT '支付方式',
    payment_status ENUM('unpaid', 'paid', 'refunded') DEFAULT 'unpaid' COMMENT '支付状态',
    status ENUM('pending', 'paid', 'completed', 'cancelled') DEFAULT 'pending' COMMENT '订单状态',
    operator_id BIGINT UNSIGNED COMMENT '操作人ID',
    remark TEXT COMMENT '备注',
    paid_at TIMESTAMP NULL COMMENT '支付时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (operator_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_order_no (order_no),
    INDEX idx_customer_phone (customer_phone),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='销售订单表';

-- =============================================
-- 销售订单明细表
-- =============================================
CREATE TABLE IF NOT EXISTS sale_order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    book_id BIGINT UNSIGNED NOT NULL COMMENT '图书ID',
    quantity INT NOT NULL COMMENT '数量',
    unit_price DECIMAL(10,2) NOT NULL COMMENT '单价',
    discount DECIMAL(5,2) DEFAULT 100 COMMENT '折扣百分比',
    total_price DECIMAL(10,2) NOT NULL COMMENT '总价',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES sale_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE RESTRICT,
    INDEX idx_order_id (order_id),
    INDEX idx_book_id (book_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='销售订单明细表';

-- =============================================
-- 购物车表
-- =============================================
CREATE TABLE IF NOT EXISTS carts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    session_id VARCHAR(100) COMMENT '会话ID',
    user_id BIGINT UNSIGNED COMMENT '用户ID',
    book_id BIGINT UNSIGNED NOT NULL COMMENT '图书ID',
    quantity INT NOT NULL DEFAULT 1 COMMENT '数量',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    INDEX idx_session_id (session_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车表';

-- =============================================
-- 操作日志表
-- =============================================
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED COMMENT '用户ID',
    username VARCHAR(50) COMMENT '用户名',
    module VARCHAR(50) COMMENT '模块',
    action VARCHAR(50) COMMENT '操作',
    target_type VARCHAR(50) COMMENT '目标类型',
    target_id BIGINT UNSIGNED COMMENT '目标ID',
    content TEXT COMMENT '操作内容',
    ip VARCHAR(45) COMMENT 'IP地址',
    user_agent VARCHAR(255) COMMENT '用户代理',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_module (module),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- =============================================
-- 系统配置表
-- =============================================
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    config_key VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
    config_value TEXT COMMENT '配置值',
    config_type VARCHAR(20) DEFAULT 'string' COMMENT '值类型',
    description VARCHAR(255) COMMENT '描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_config_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- =============================================
-- 初始化数据
-- =============================================

-- 插入默认管理员 (密码: admin123)
INSERT INTO users (username, password, email, real_name, status) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'admin@library.com', '系统管理员', 1);

-- 插入默认角色
INSERT INTO roles (name, code, description, sort_order) VALUES
('超级管理员', 'super_admin', '拥有系统所有权限', 1),
('管理员', 'admin', '系统管理员，拥有大部分权限', 2),
('图书管理员', 'librarian', '管理图书和借阅', 3),
('普通用户', 'user', '普通用户，基本权限', 4);

-- 插入默认权限
INSERT INTO permissions (name, code, resource, action, description, sort_order) VALUES
-- 用户管理
('用户管理', 'user', 'user', '*', '用户管理权限', 1),
('用户查看', 'user:read', 'user', 'read', '查看用户列表', 2),
('用户创建', 'user:create', 'user', 'create', '创建用户', 3),
('用户编辑', 'user:update', 'user', 'update', '编辑用户', 4),
('用户删除', 'user:delete', 'user', 'delete', '删除用户', 5),
-- 角色管理
('角色管理', 'role', 'role', '*', '角色管理权限', 10),
('角色查看', 'role:read', 'role', 'read', '查看角色列表', 11),
-- 图书管理
('图书管理', 'book', 'book', '*', '图书管理权限', 20),
('图书查看', 'book:read', 'book', 'read', '查看图书列表', 21),
('图书创建', 'book:create', 'book', 'create', '创建图书', 22),
('图书编辑', 'book:update', 'book', 'update', '编辑图书', 23),
('图书删除', 'book:delete', 'book', 'delete', '删除图书', 24),
-- 库存管理
('库存管理', 'stock', 'stock', '*', '库存管理权限', 30),
('库存查看', 'stock:read', 'stock', 'read', '查看库存', 31),
('入库操作', 'stock:in', 'stock', 'in', '图书入库', 32),
('出库操作', 'stock:out', 'stock', 'out', '图书出库', 33),
-- 采购管理
('采购管理', 'purchase', 'purchase', '*', '采购管理权限', 40),
('采购查看', 'purchase:read', 'purchase', 'read', '查看采购订单', 41),
('采购创建', 'purchase:create', 'purchase', 'create', '创建采购订单', 42),
-- 销售管理
('销售管理', 'sale', 'sale', '*', '销售管理权限', 50),
('销售查看', 'sale:read', 'sale', 'read', '查看销售订单', 51),
('销售创建', 'sale:create', 'sale', 'create', '创建销售订单', 52),
-- 借阅管理
('借阅管理', 'borrow', 'borrow', '*', '借阅管理权限', 60),
('借阅查看', 'borrow:read', 'borrow', 'read', '查看借阅记录', 61),
('借阅操作', 'borrow:operate', 'borrow', 'operate', '借阅操作', 62);

-- 分配管理员角色
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);

-- 分配超级管理员所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- 分配管理员部分权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE code IN ('user:read', 'book', 'stock', 'purchase', 'sale', 'borrow');

-- 插入系统配置
INSERT INTO system_configs (config_key, config_value, config_type, description) VALUES
('site_name', '图书管理系统', 'string', '网站名称'),
('borrow_days', '30', 'number', '默认借阅天数'),
('max_borrow_books', '5', 'number', '最大借阅数量'),
('max_renew_times', '2', 'number', '最大续借次数'),
('fine_per_day', '0.50', 'number', '逾期罚款(元/天)'),
('low_stock_threshold', '5', 'number', '库存预警阈值');

-- 插入测试图书数据
INSERT INTO books (isbn, title, author, publisher, publish_date, category, price, description, pages, language) VALUES
('978-7-111-54784-2', '深入理解计算机系统', 'Randal E. Bryant', '机械工业出版社', '2016-07-01', '计算机', 139.00, '本书从程序员的视角详细阐述计算机系统的本质概念', 737, '中文'),
('978-7-115-42533-4', 'JavaScript高级程序设计', 'Matt Frisbie', '人民邮电出版社', '2020-05-01', '计算机', 129.00, 'JavaScript经典著作，前端开发必读', 912, '中文'),
('978-7-115-42602-7', 'Python编程从入门到实践', 'Eric Matthes', '人民邮电出版社', '2020-10-01', '计算机', 89.00, 'Python入门经典教程', 459, '中文'),
('978-7-111-40701-0', '算法导论', 'Thomas H. Cormen', '机械工业出版社', '2013-01-01', '计算机', 128.00, '算法领域的经典著作', 780, '中文'),
('978-7-115-42813-7', '设计模式', 'Erich Gamma', '人民邮电出版社', '2010-07-01', '计算机', 99.00, '软件设计模式经典', 395, '中文');

-- 为图书创建库存记录
INSERT INTO stocks (book_id, total_quantity, available_quantity, location, warning_threshold)
SELECT id, 20, 20, CONCAT('A-', LPAD(id, 2, '0'), '-01'), 5 FROM books;

-- 插入测试供应商
INSERT INTO suppliers (name, code, contact, phone, email, address, status) VALUES
('北京图书批发中心', 'BJ-001', '张经理', '010-12345678', 'bj@supplier.com', '北京市朝阳区图书大厦', 1),
('上海书城', 'SH-001', '李经理', '021-87654321', 'sh@supplier.com', '上海市黄浦区书城大厦', 1),
('广州图书市场', 'GZ-001', '王经理', '020-11112222', 'gz@supplier.com', '广州市天河区图书市场', 1);

-- 完成
SELECT 'Database initialization completed!' AS message;
