'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { 
  Card, 
  Button, 
  Table, 
  Modal, 
  Form, 
  Input, 
  Select, 
  DatePicker, 
  message, 
  Space, 
  Tag,
  Popconfirm 
} from 'antd';
import { 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  LogoutOutlined,
  CheckOutlined,
  UndoOutlined
} from '@ant-design/icons';
import { apiClient } from '@/lib/api';
import { useAuth } from '@/lib/auth';
import { Task, CreateTaskRequest, UpdateTaskRequest } from '@/types';
import dayjs from 'dayjs';

const { TextArea } = Input;
const { Option } = Select;

export default function TasksPage() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [showCompleted, setShowCompleted] = useState(false); // 完了タスクの表示制御
  const [form] = Form.useForm();
  const router = useRouter();
  const { user, logout, loading: authLoading } = useAuth();

  // タスク一覧を取得
  const fetchTasks = async () => {
    setLoading(true);
    try {
      const response = await apiClient.getTasks();
      if (response.success) {
        // データがnullの場合は空配列を設定
        setTasks(response.data || []);
      }
    } catch (error) {
      message.error('タスクの取得に失敗しました');
      console.error('Fetch tasks error:', error);
      setTasks([]); // エラー時も空配列を設定
    } finally {
      setLoading(false);
    }
  };

  // タスクを作成
  const handleCreateTask = async (values: any) => {
    try {
      const taskData: CreateTaskRequest = {
        title: values.title,
        description: values.description || undefined,
        deadline: values.deadline ? values.deadline.toDate() : undefined,
        priority: values.priority,
      };

      const response = await apiClient.createTask(taskData);
      if (response.success && response.data) {
        message.success('タスクを作成しました');
        setModalVisible(false);
        form.resetFields();
        fetchTasks();
      }
    } catch (error) {
      message.error('タスクの作成に失敗しました');
      console.error('Create task error:', error);
    }
  };

  // タスクを更新
  const handleUpdateTask = async (values: any) => {
    if (!editingTask) return;

    try {
      const updateData: UpdateTaskRequest = {
        title: values.title,
        description: values.description || undefined,
        deadline: values.deadline ? values.deadline.toDate() : undefined,
        priority: values.priority,
        status: values.status,
      };

      const response = await apiClient.updateTask(editingTask.id, updateData);
      if (response.success && response.data) {
        message.success('タスクを更新しました');
        setModalVisible(false);
        setEditingTask(null);
        form.resetFields();
        fetchTasks();
      }
    } catch (error) {
      message.error('タスクの更新に失敗しました');
      console.error('Update task error:', error);
    }
  };

  // タスクを削除
  const handleDeleteTask = async (taskId: string) => {
    try {
      const response = await apiClient.deleteTask(taskId);
      if (response.success) {
        message.success('タスクを削除しました');
        fetchTasks();
      }
    } catch (error) {
      message.error('タスクの削除に失敗しました');
      console.error('Delete task error:', error);
    }
  };

  // タスクのステータスを切り替え
  const handleToggleTaskStatus = async (task: Task) => {
    try {
      const newStatus = task.status === 'pending' ? 'completed' : 'pending';
      const updateData: UpdateTaskRequest = {
        title: task.title,
        description: task.description,
        deadline: task.deadline,
        priority: task.priority,
        status: newStatus,
      };

      const response = await apiClient.updateTask(task.id, updateData);
      if (response.success) {
        message.success(`タスクを${newStatus === 'completed' ? '完了' : '未完了'}に変更しました`);
        fetchTasks();
      }
    } catch (error) {
      message.error('タスクのステータス変更に失敗しました');
      console.error('Toggle task status error:', error);
    }
  };

  // モーダルを開く（新規作成）
  const openCreateModal = () => {
    setEditingTask(null);
    setModalVisible(true);
    form.resetFields();
  };

  // モーダルを開く（編集）
  const openEditModal = (task: Task) => {
    setEditingTask(task);
    setModalVisible(true);
    form.setFieldsValue({
      title: task.title,
      description: task.description,
      deadline: task.deadline ? dayjs(task.deadline) : undefined,
      priority: task.priority,
      status: task.status,
    });
  };

  // ログアウト
  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  // タスクを未完了と完了済みに分ける
  const pendingTasks = tasks.filter(task => task.status === 'pending');
  const completedTasks = tasks.filter(task => task.status === 'completed');

  // 優先度の色を取得
  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return 'red';
      case 'medium': return 'orange';
      case 'low': return 'green';
      default: return 'default';
    }
  };

  // ステータスの色を取得
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return 'green';
      case 'pending': return 'blue';
      default: return 'default';
    }
  };

  // テーブルの列定義
  const columns = [
    {
      title: 'タイトル',
      dataIndex: 'title',
      key: 'title',
    },
    {
      title: '説明',
      dataIndex: 'description',
      key: 'description',
      render: (text: string) => text || '-',
    },
    {
      title: '期限',
      dataIndex: 'deadline',
      key: 'deadline',
      render: (date: string) => date ? dayjs(date).format('YYYY-MM-DD') : '-',
    },
    {
      title: '優先度',
      dataIndex: 'priority',
      key: 'priority',
      render: (priority: string) => (
        <Tag color={getPriorityColor(priority)}>
          {priority === 'high' ? '高' : priority === 'medium' ? '中' : '低'}
        </Tag>
      ),
    },
    {
      title: 'ステータス',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>
          {status === 'completed' ? '完了' : '未完了'}
        </Tag>
      ),
    },
    {
      title: '作成日',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record: Task) => (
        <Space>
          <Button
            type="link"
            icon={record.status === 'pending' ? <CheckOutlined /> : <UndoOutlined />}
            onClick={() => handleToggleTaskStatus(record)}
          >
            {record.status === 'pending' ? '完了' : '未完了に戻す'}
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => openEditModal(record)}
          >
            編集
          </Button>
          <Popconfirm
            title="このタスクを削除しますか？"
            onConfirm={() => handleDeleteTask(record.id)}
            okText="削除"
            cancelText="キャンセル"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
            >
              削除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  useEffect(() => {
    // 認証のローディングが完了してから処理
    if (authLoading) {
      return;
    }
    
    if (!user) {
      router.push('/login');
      return;
    }
    
    fetchTasks();
  }, [user, router, authLoading]);

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

  // ユーザーが認証されていない場合は何も表示しない（リダイレクト中）
  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* ヘッダー */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">タスク管理</h1>
            <p className="text-gray-600">こんにちは、{user.name}さん</p>
          </div>
           <div className="flex space-x-4">
             <Button
               type="primary"
               icon={<PlusOutlined />}
               onClick={openCreateModal}
             >
               新しいタスク
             </Button>
             <Button
               type={showCompleted ? "default" : "primary"}
               onClick={() => setShowCompleted(!showCompleted)}
             >
               {showCompleted ? '未完了タスクを表示' : '完了済みタスクを表示'}
             </Button>
             <Button
               icon={<LogoutOutlined />}
               onClick={handleLogout}
             >
               ログアウト
             </Button>
           </div>
        </div>

        {/* タスク一覧 */}
        <div className="space-y-6">
          {/* 未完了タスク */}
          <Card>
            <div className="mb-4">
              <h2 className="text-xl font-semibold text-gray-800">
                未完了タスク ({pendingTasks.length}件)
              </h2>
            </div>
            <Table
              columns={columns}
              dataSource={pendingTasks}
              rowKey="id"
              loading={loading}
              pagination={{
                pageSize: 10,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total, range) =>
                  `${range[0]}-${range[1]} / ${total}件`,
              }}
            />
          </Card>

          {/* 完了済みタスク */}
          {completedTasks.length > 0 && (
            <Card className="bg-gray-50">
              <div className="mb-4 flex justify-between items-center">
                <h2 className="text-xl font-semibold text-gray-600">
                  完了済みタスク ({completedTasks.length}件)
                </h2>
                <Button
                  type="link"
                  onClick={() => setShowCompleted(!showCompleted)}
                >
                  {showCompleted ? '非表示' : '表示'}
                </Button>
              </div>
              {showCompleted && (
                <Table
                  columns={columns}
                  dataSource={completedTasks}
                  rowKey="id"
                  loading={loading}
                  pagination={{
                    pageSize: 10,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (total, range) =>
                      `${range[0]}-${range[1]} / ${total}件`,
                  }}
                />
              )}
            </Card>
          )}
        </div>

        {/* タスク作成/編集モーダル */}
        <Modal
          title={editingTask ? 'タスクを編集' : '新しいタスク'}
          open={modalVisible}
          onCancel={() => {
            setModalVisible(false);
            setEditingTask(null);
            form.resetFields();
          }}
          footer={null}
          width={600}
        >
          <Form
            form={form}
            layout="vertical"
            onFinish={editingTask ? handleUpdateTask : handleCreateTask}
          >
            <Form.Item
              name="title"
              label="タイトル"
              rules={[{ required: true, message: 'タイトルを入力してください' }]}
            >
              <Input placeholder="タスクのタイトル" />
            </Form.Item>

            <Form.Item
              name="description"
              label="説明"
            >
              <TextArea
                rows={4}
                placeholder="タスクの詳細説明（任意）"
              />
            </Form.Item>

            <Form.Item
              name="deadline"
              label="期限"
            >
              <DatePicker
                style={{ width: '100%' }}
                placeholder="期限を選択（任意）"
                format="YYYY-MM-DD"
              />
            </Form.Item>

            <Form.Item
              name="priority"
              label="優先度"
              rules={[{ required: true, message: '優先度を選択してください' }]}
            >
              <Select placeholder="優先度を選択">
                <Option value="high">高</Option>
                <Option value="medium">中</Option>
                <Option value="low">低</Option>
              </Select>
            </Form.Item>

            {editingTask && (
              <Form.Item
                name="status"
                label="ステータス"
                rules={[{ required: true, message: 'ステータスを選択してください' }]}
              >
                <Select placeholder="ステータスを選択">
                  <Option value="pending">未完了</Option>
                  <Option value="completed">完了</Option>
                </Select>
              </Form.Item>
            )}

            <Form.Item className="mb-0">
              <Space>
                <Button type="primary" htmlType="submit">
                  {editingTask ? '更新' : '作成'}
                </Button>
                <Button onClick={() => {
                  setModalVisible(false);
                  setEditingTask(null);
                  form.resetFields();
                }}>
                  キャンセル
                </Button>
              </Space>
            </Form.Item>
          </Form>
        </Modal>
      </div>
    </div>
  );
}
