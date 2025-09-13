'use client';

import { useEffect } from 'react';

export function AntdProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    // React 19互換性パッチを確実に適用
    import('@ant-design/v5-patch-for-react-19');
  }, []);

  return <>{children}</>;
}

