'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Form, Input, Button, Card, message } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { apiClient } from '@/lib/api';
import { useAuth } from '@/lib/auth';

export default function RegisterPage() {
  const [loading, setLoading] = useState(false);
  const router = useRouter();
  const { login, user, loading: authLoading } = useAuth();

  // 既にログインしている場合はタスクページにリダイレクト
  useEffect(() => {
    if (!authLoading && user) {
      router.push('/tasks');
    }
  }, [user, authLoading, router]);

  const onFinish = async (values: { email: string; password: string; name: string }) => {
    setLoading(true);
    try {
      const response = await apiClient.register(values);
      if (response.success && response.data) {
        login(response.data.access_token, response.data.user);
        message.success('アカウントを作成しました');
        router.push('/tasks');
      }
    } catch (error) {
      message.error('アカウント作成に失敗しました');
      console.error('Register error:', error);
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
            新規アカウント作成
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            または{' '}
            <a
              href="/login"
              className="font-medium text-blue-600 hover:text-blue-500"
            >
              既存のアカウントでログイン
            </a>
          </p>
        </div>
        
        <Card className="shadow-lg">
          <Form
            name="register"
            onFinish={onFinish}
            autoComplete="off"
            layout="vertical"
          >
            <Form.Item
              name="name"
              label="ユーザー名"
              rules={[
                { required: true, message: 'ユーザー名を入力してください' },
                { min: 2, message: 'ユーザー名は2文字以上で入力してください' }
              ]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="ユーザー名"
                size="large"
              />
            </Form.Item>

            <Form.Item
              name="email"
              label="メールアドレス"
              rules={[
                { required: true, message: 'メールアドレスを入力してください' },
                { type: 'email', message: '有効なメールアドレスを入力してください' }
              ]}
            >
              <Input
                prefix={<MailOutlined />}
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

            <Form.Item
              name="confirmPassword"
              label="パスワード確認"
              dependencies={['password']}
              rules={[
                { required: true, message: 'パスワードを再入力してください' },
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    if (!value || getFieldValue('password') === value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error('パスワードが一致しません'));
                  },
                }),
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="パスワード再入力"
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
                アカウント作成
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </div>
    </div>
  );
}
