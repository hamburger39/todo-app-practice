// ユーザー関連の型定義
export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

// タスク関連の型定義
export interface Task {
  id: string;
  user_id: string;
  title: string;
  description?: string;
  deadline?: string;
  priority: 'high' | 'medium' | 'low';
  status: 'pending' | 'completed';
  created_at: string;
  updated_at: string;
}

export interface CreateTaskRequest {
  title: string;
  description?: string;
  deadline?: Date;
  priority: 'high' | 'medium' | 'low';
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  deadline?: Date;
  priority?: 'high' | 'medium' | 'low';
  status?: 'pending' | 'completed';
}

// APIレスポンスの型定義
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

// フィルター・ソート関連
export interface TaskFilters {
  status?: 'pending' | 'completed' | 'all';
  priority?: 'high' | 'medium' | 'low' | 'all';
  sortBy?: 'deadline' | 'priority' | 'created_at';
  sortOrder?: 'asc' | 'desc';
}
