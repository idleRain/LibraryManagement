// 用户相关类型
export interface User {
  id: number;
  username: string;
  email: string;
  phone?: string;
  real_name?: string;
  status: number;
  roles: Role[];
  created_at: string;
  updated_at: string;
}

export interface Role {
  id: number;
  name: string;
  code: string;
  description?: string;
  permissions?: Permission[];
  created_at: string;
  updated_at: string;
}

export interface Permission {
  id: number;
  name: string;
  code: string;
  resource: string;
  action: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

// 图书相关类型
export interface Book {
  id: number;
  isbn: string;
  title: string;
  author?: string;
  publisher?: string;
  publish_date?: string;
  category?: string;
  price: number;
  description?: string;
  cover_image?: string;
  pages?: number;
  language?: string;
  stock?: Stock;
  created_at: string;
  updated_at: string;
}

export interface Stock {
  id: number;
  book_id: number;
  book?: Book;
  total_quantity: number;
  available_quantity: number;
  borrowed_quantity: number;
  sold_quantity: number;
  location?: string;
  created_at: string;
  updated_at: string;
}

export interface StockRecord {
  id: number;
  book_id: number;
  book?: Book;
  type: 'in' | 'out';
  quantity: number;
  reason?: string;
  operator_id?: number;
  operator?: User;
  created_at: string;
}

// 采购相关类型
export interface Supplier {
  id: number;
  name: string;
  contact?: string;
  phone?: string;
  email?: string;
  address?: string;
  description?: string;
  status: number;
  created_at: string;
  updated_at: string;
}

export interface PurchaseOrder {
  id: number;
  order_no: string;
  supplier_id: number;
  supplier?: Supplier;
  total_amount: number;
  status: 'pending' | 'approved' | 'completed' | 'cancelled';
  operator_id?: number;
  operator?: User;
  remark?: string;
  items: PurchaseOrderItem[];
  created_at: string;
  updated_at: string;
}

export interface PurchaseOrderItem {
  id: number;
  order_id: number;
  book_id: number;
  book?: Book;
  quantity: number;
  unit_price: number;
  total_price: number;
  received_quantity: number;
  created_at: string;
  updated_at: string;
}

// 销售相关类型
export interface SaleOrder {
  id: number;
  order_no: string;
  customer_name?: string;
  customer_phone?: string;
  total_amount: number;
  pay_amount: number;
  discount: number;
  status: 'pending' | 'paid' | 'completed' | 'cancelled';
  payment_method?: string;
  operator_id?: number;
  operator?: User;
  remark?: string;
  items: SaleOrderItem[];
  created_at: string;
  updated_at: string;
}

export interface SaleOrderItem {
  id: number;
  order_id: number;
  book_id: number;
  book?: Book;
  quantity: number;
  unit_price: number;
  total_price: number;
  discount: number;
  created_at: string;
  updated_at: string;
}

// 借阅相关类型
export interface BorrowRecord {
  id: string;
  book_id: number;
  book_title: string;
  book_isbn: string;
  user_id: number;
  user_name: string;
  borrow_date: string;
  due_date: string;
  return_date?: string;
  status: 'borrowed' | 'returned' | 'overdue' | 'lost';
  renew_count: number;
  fine: number;
  fine_paid: boolean;
  operator_id?: number;
  remark?: string;
  created_at: string;
  updated_at: string;
}

// API 响应类型
export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data?: T;
}

export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

// 登录响应
export interface LoginResponse {
  token: string;
  user: User;
  permissions: string[];
}

// 统计数据
export interface DashboardStats {
  total_books: number;
  total_users: number;
  total_borrowed: number;
  total_sales: number;
  recent_borrows: BorrowRecord[];
  low_stock_books: Book[];
}

export interface SalesStats {
  total_orders: number;
  total_amount: number;
  total_quantity: number;
}

export interface BorrowStats {
  current_borrowed: number;
  total_borrowed: number;
  overdue_count: number;
}
