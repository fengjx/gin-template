import { LoadingScreen } from '@/components/shared/loading-screen';
import { Alert, AlertDescription, AlertIcon, AlertTitle } from '@/components/ui/alert';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { refreshSession } from '../api/client';
import { useAuth } from '../features/auth/AuthContext';

export function OAuthCallbackPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { reload } = useAuth();
  const [error, setError] = useState('');

  useEffect(() => {
    const status = searchParams.get('status');
    if (status !== 'success') {
      setError(searchParams.get('message') || 'OAuth 登录失败');
      return;
    }
    refreshSession()
      .then(async () => {
        await reload();
        navigate('/');
      })
      .catch((err) => setError((err as Error).message));
  }, [navigate, reload, searchParams]);

  return (
    <div className="grid min-h-screen place-items-center px-4">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>OAuth 登录回调</CardTitle>
        </CardHeader>
        <CardContent>
          {error ? (
            <Alert variant="destructive">
              <AlertIcon className="h-4 w-4" />
              <AlertTitle>登录失败</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          ) : (
            <LoadingScreen message="正在完成 OAuth 登录..." />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
