'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Form, Input, Button, Card, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { apiClient } from '@/lib/api';
import { useAuth } from '@/lib/auth';

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const router = useRouter();
  const { login, user, loading: authLoading } = useAuth();

  // 既にログインしている場合はタスクページにリダイレクト
  useEffect(() => {
    if (!authLoading && user) {
      router.push('/tasks');
    }
  }, [user, authLoading, router]);

  const onFinish = async (values: { email: string; password: string }) => {
    setLoading(true);
    try {
      const response = await apiClient.login(values);
      if (response.success && response.data) {
        login(response.data.access_token, response.data.user);
        message.success('ログインに成功しました');
        router.push('/tasks');
      }
    } catch (error) {
      message.error('ログインに失敗しました');
      console.error('Login error:', error);
    } finally {
      setLoading(false);
    }
  };

  // 認証のローディング中は何も表示しない
  if (authLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">読み込み中...</p>
        </div>
      </div>
    );
  }

  // 既にログインしている場合は何も表示しない（リダイレクト中）
  if (user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            アカウントにログイン
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            または{' '}
            <a
              href="/register"
              className="font-medium text-blue-600 hover:text-blue-500"
            >
              新規登録
            </a>
          </p>
        </div>
        
        <Card className="shadow-lg">
          <Form
            name="login"
            onFinish={onFinish}
            autoComplete="off"
            layout="vertical"
          >
            <Form.Item
              name="email"
              label="メールアドレス"
              rules={[
                { required: true, message: 'メールアドレスを入力してください' },
                { type: 'email', message: '有効なメールアドレスを入力してください' }
              ]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="example@email.com"
                size="large"
              />
            </Form.Item>

            <Form.Item
              name="password"
              label="パスワード"
              rules={[
                { required: true, message: 'パスワードを入力してください' },
                { min: 6, message: 'パスワードは6文字以上で入力してください' }
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="パスワード"
                size="large"
              />
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                size="large"
                className="w-full"
              >
                ログイン
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </div>
    </div>
  );
}
