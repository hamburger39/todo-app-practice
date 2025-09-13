'use client';

import { useState } from 'react';
import { Button, Card, message } from 'antd';
import { apiClient } from '@/lib/api';

export default function TestPage() {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<string>('');

  const testHealthCheck = async () => {
    setLoading(true);
    try {
      // 直接fetchでテスト
      const response = await fetch('http://localhost:8080/health');
      const data = await response.json();
      setResult(JSON.stringify(data, null, 2));
      message.success('Health check successful');
    } catch (error) {
      setResult(`Error: ${error instanceof Error ? error.message : String(error)}`);
      message.error('Health check failed');
    } finally {
      setLoading(false);
    }
  };

  const testRegister = async () => {
    setLoading(true);
    try {
      const response = await apiClient.register({
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User'
      });
      setResult(JSON.stringify(response, null, 2));
      message.success('Register test successful');
    } catch (error) {
      setResult(`Error: ${error instanceof Error ? error.message : String(error)}`);
      message.error('Register test failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">API Test Page</h1>
        
        <div className="space-y-4">
          <Card title="Health Check Test">
            <Button 
              onClick={testHealthCheck} 
              loading={loading}
              type="primary"
            >
              Test Health Check
            </Button>
          </Card>

          <Card title="Register Test">
            <Button 
              onClick={testRegister} 
              loading={loading}
              type="primary"
            >
              Test Register API
            </Button>
          </Card>

          {result && (
            <Card title="Result">
              <pre className="bg-gray-100 p-4 rounded overflow-auto">
                {result}
              </pre>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}
