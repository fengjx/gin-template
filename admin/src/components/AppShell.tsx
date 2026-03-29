import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import { cn } from '@/lib/utils';
import {
  Activity,
  FileText,
  LayoutDashboard,
  LogOut,
  Menu,
  Settings2,
  Sparkles,
  UserCircle2,
  Users,
  Workflow,
} from 'lucide-react';
import { useState } from 'react';
import { Link, NavLink, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../features/auth/AuthContext';

type MenuItem = {
  path: string;
  label: string;
  description: string;
  icon: typeof LayoutDashboard;
};

const menuItems: MenuItem[] = [
  { path: '/', label: '工作台', description: '统一查看系统看板', icon: LayoutDashboard },
  { path: '/users', label: '用户管理', description: '账号状态、角色与邮箱验证', icon: Users },
  { path: '/files', label: '文件管理', description: '上传资产、容量与清理操作', icon: FileText },
  {
    path: '/options',
    label: '系统配置',
    description: '维护系统级配置项与公开属性',
    icon: Settings2,
  },
  {
    path: '/profile',
    label: '个人中心',
    description: '维护当前账号资料与展示信息',
    icon: UserCircle2,
  },
  { path: '/docs', label: '接口文档', description: 'OpenAPI 文档与联调入口', icon: Workflow },
  { path: '/pprof', label: '性能监控', description: 'Pprof 页面与性能分析入口', icon: Activity },
];

function getInitial(value: string) {
  return value.trim().charAt(0).toUpperCase() || 'G';
}

function matchMenuItem(pathname: string) {
  if (pathname === '/') {
    return menuItems[0];
  }
  return (
    menuItems.find((item) => item.path !== '/' && pathname.startsWith(item.path)) ?? menuItems[0]
  );
}

type SidebarProps = {
  onNavigate?: () => void;
};

function SidebarContent({ onNavigate }: SidebarProps) {
  return (
    <div className="flex h-full flex-col">
      <Link
        to="/"
        className="flex items-center gap-4 border-b border-white/10 px-6 py-6"
        onClick={onNavigate}
      >
        <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-teal-400 to-sky-400 text-sm font-black tracking-[0.28em] text-slate-950 shadow-glow">
          GT
        </div>
        <div className="space-y-1">
          <div className="font-display text-lg font-semibold tracking-wide text-white">
            gin-template
          </div>
          <div className="text-xs uppercase tracking-[0.24em] text-slate-400">
            Go · React · OpenAPI
          </div>
        </div>
      </Link>

      <div className="space-y-2 px-4 py-6">
        {menuItems.map((item) => {
          const Icon = item.icon;
          return (
            <NavLink
              key={item.path}
              to={item.path}
              onClick={onNavigate}
              className={({ isActive }) =>
                cn(
                  'group flex items-center gap-4 rounded-[1.4rem] px-4 py-3 transition-all',
                  isActive
                    ? 'bg-white/10 text-white shadow-[inset_0_0_0_1px_rgba(255,255,255,0.08)]'
                    : 'text-slate-300 hover:bg-white/6 hover:text-white',
                )
              }
            >
              {({ isActive }) => (
                <>
                  <div
                    className={cn(
                      'flex h-10 w-10 items-center justify-center rounded-2xl border transition-colors',
                      isActive
                        ? 'border-white/10 bg-white/10 text-white'
                        : 'border-white/5 bg-white/5 text-slate-300',
                    )}
                  >
                    <Icon className="h-5 w-5" />
                  </div>
                  <div className="min-w-0">
                    <div className="font-medium">{item.label}</div>
                    <div className="truncate text-xs text-slate-400">{item.description}</div>
                  </div>
                </>
              )}
            </NavLink>
          );
        })}
      </div>
    </div>
  );
}

export function AppShell() {
  const [mobileOpen, setMobileOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { logout, user } = useAuth();
  const displayName = user?.display_name || user?.username || '未登录用户';
  const currentItem = matchMenuItem(location.pathname);

  async function handleLogout() {
    await logout();
    navigate('/login');
  }

  return (
    <div className="relative min-h-screen overflow-hidden">
      <div className="pointer-events-none absolute inset-0 page-dots opacity-60" />
      <div className="pointer-events-none absolute inset-x-0 top-0 h-[34rem] bg-gradient-to-b from-teal-100/70 via-sky-50/50 to-transparent" />

      <div className="relative flex min-h-screen">
        <aside className="hidden w-[296px] shrink-0 px-6 py-6 lg:block">
          <div className="sticky top-6 h-[calc(100vh-3rem)] overflow-hidden rounded-[2rem] border border-white/10 bg-sidebar text-sidebar-foreground shadow-[0_28px_90px_rgba(2,6,23,0.28)]">
            <SidebarContent />
          </div>
        </aside>

        <div className="flex min-h-screen flex-1 flex-col">
          <header className="px-4 pt-4 sm:px-6 lg:px-8">
            <div className="glass-panel flex items-center justify-between gap-4 rounded-[1.75rem] px-4 py-3 sm:px-6">
              <div className="flex min-w-0 items-center gap-3">
                <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
                  <SheetTrigger asChild>
                    <Button variant="outline" size="icon" className="lg:hidden">
                      <Menu className="h-5 w-5" />
                      <span className="sr-only">打开导航</span>
                    </Button>
                  </SheetTrigger>
                  <SheetContent
                    side="left"
                    className="w-[320px] border-0 bg-sidebar p-0 text-sidebar-foreground"
                  >
                    <SheetHeader className="sr-only">
                      <SheetTitle>导航菜单</SheetTitle>
                      <SheetDescription>在移动端切换后台页面。</SheetDescription>
                    </SheetHeader>
                    <SidebarContent onNavigate={() => setMobileOpen(false)} />
                  </SheetContent>
                </Sheet>

                <div className="min-w-0">
                  <div className="text-xs font-semibold uppercase tracking-[0.24em] text-primary">
                    Console
                  </div>
                  <div className="truncate text-sm text-muted-foreground sm:text-base">
                    {currentItem.label} · {currentItem.description}
                  </div>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <button
                      type="button"
                      className="flex items-center gap-3 rounded-full border border-border/60 bg-background/80 px-2 py-2 pl-2 pr-4 text-left shadow-sm transition-colors hover:bg-accent focus:outline-none focus:ring-2 focus:ring-ring"
                    >
                      <Avatar className="h-10 w-10 border border-primary/10">
                        <AvatarFallback>{getInitial(displayName)}</AvatarFallback>
                      </Avatar>
                      <div className="hidden min-w-0 sm:block">
                        <div className="truncate text-sm font-semibold text-foreground">
                          {displayName}
                        </div>
                        <div className="truncate text-xs text-muted-foreground">
                          {user?.email || '当前登录用户'}
                        </div>
                      </div>
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-64">
                    <DropdownMenuLabel>账号中心</DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onSelect={() => navigate('/profile')}>
                      <UserCircle2 className="h-4 w-4" />
                      个人中心
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      className="text-rose-600 focus:text-rose-600"
                      onSelect={() => void handleLogout()}
                    >
                      <LogOut className="h-4 w-4" />
                      退出登录
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          </header>

          <main className="flex-1 px-4 pb-8 pt-6 sm:px-6 lg:px-8">
            <div className="mx-auto w-full max-w-[1480px]">
              <Outlet />
            </div>
          </main>
        </div>
      </div>
    </div>
  );
}
