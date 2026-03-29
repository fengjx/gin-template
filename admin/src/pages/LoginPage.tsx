import { Alert, AlertDescription, AlertIcon, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { KeyRound, ShieldCheck } from 'lucide-react';
import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../features/auth/AuthContext';

export function LoginPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { login } = useAuth();
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [loginForm, setLoginForm] = useState({
    identifier: '',
    password: '',
  });

  async function handleLogin(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      setSubmitting(true);
      setError('');
      await login(loginForm.identifier, loginForm.password);
      const redirect = searchParams.get('redirect');
      navigate(redirect?.startsWith('/') ? redirect : '/');
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="relative min-h-screen overflow-hidden bg-slate-950">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(45,212,191,0.16),transparent_24%),radial-gradient(circle_at_bottom_right,rgba(56,189,248,0.14),transparent_28%),linear-gradient(135deg,#020617_0%,#0f172a_52%,#0f3f46_100%)]" />
      <div className="absolute inset-0 page-dots opacity-15" />

      <div className="relative flex min-h-screen items-center justify-center px-6 py-10">
        <Card className="w-full max-w-[460px] border-white/20 bg-white text-slate-950 shadow-[0_28px_80px_rgba(2,6,23,0.42)]">
          <CardHeader className="space-y-4 pb-4 text-center">
            <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br from-teal-400 to-sky-400 text-base font-black tracking-[0.28em] text-slate-950">
              GT
            </div>
            <div className="inline-flex w-fit items-center gap-2 self-center rounded-full border border-teal-200 bg-teal-50 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-teal-700">
              <ShieldCheck className="h-3.5 w-3.5" />
              Secure Access
            </div>
            <div className="space-y-2">
              <CardTitle className="font-display text-3xl font-semibold tracking-tight text-slate-950">
                登录
              </CardTitle>
              <CardDescription className="text-sm font-medium leading-7 text-slate-700">
                使用你的账号进入后台工作台。
              </CardDescription>
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            {error ? (
              <Alert variant="destructive">
                <AlertIcon className="h-4 w-4" />
                <AlertTitle>登录失败</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            ) : null}

            <form className="space-y-5" onSubmit={handleLogin}>
              <div className="space-y-2">
                <Label htmlFor="login-identifier" className="text-sm font-semibold text-slate-900">
                  账号
                </Label>
                <Input
                  id="login-identifier"
                  placeholder="用户名或邮箱"
                  autoComplete="username"
                  value={loginForm.identifier}
                  onChange={(event) =>
                    setLoginForm((value) => ({ ...value, identifier: event.target.value }))
                  }
                  className="border-slate-300 bg-white text-slate-950 placeholder:text-slate-500"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="login-password" className="text-sm font-semibold text-slate-900">
                  密码
                </Label>
                <Input
                  id="login-password"
                  type="password"
                  autoComplete="current-password"
                  value={loginForm.password}
                  onChange={(event) =>
                    setLoginForm((value) => ({ ...value, password: event.target.value }))
                  }
                  className="border-slate-300 bg-white text-slate-950 placeholder:text-slate-500"
                  required
                />
              </div>
              <Button className="h-12 w-full text-base" type="submit" disabled={submitting}>
                <KeyRound className="h-4 w-4" />
                {submitting ? '登录中...' : '登录进入工作台'}
              </Button>
            </form>

            <div className="rounded-[1.2rem] border border-slate-200 bg-slate-50 px-4 py-3 text-left">
              <div className="text-xs font-semibold uppercase tracking-[0.2em] text-slate-700">
                Tips
              </div>
              <p className="mt-2 text-sm font-medium leading-6 text-slate-800">
                请输入用户名或邮箱进行登录。如果当前账号无法访问后台，请先确认权限或联系管理员处理。
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
