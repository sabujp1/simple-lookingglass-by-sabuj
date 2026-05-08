'use client';

import { useEffect } from 'react';
import Link from 'next/link';
import { 
  Activity, 
  Globe, 
  Server, 
  Users, 
  TrendingUp, 
  ArrowUpRight,
  ArrowDownRight,
  Clock
} from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Progress } from '@/components/ui/progress';
import { useQuery } from '@tanstack/react-query';
import { adminApi } from '@/lib/api';
import { formatDateTime } from '@/lib/utils';

const stats = [
  { name: 'Total Queries', value: '12,847', change: '+12.5%', trend: 'up', icon: Activity },
  { name: 'Active Routers', value: '24', change: '+2', trend: 'up', icon: Server },
  { name: 'Online Users', value: '156', change: '+8.2%', trend: 'up', icon: Users },
  { name: 'Avg Response Time', value: '142ms', change: '-5.3%', trend: 'down', icon: Clock },
];

const recentQueries = [
  { id: 'q1', type: 'ping', target: '8.8.8.8', router: 'core-rtr-01', status: 'completed', time: '2 min ago' },
  { id: 'q2', type: 'traceroute', target: 'google.com', router: 'edge-rtr-02', status: 'completed', time: '5 min ago' },
  { id: 'q3', type: 'bgp', target: 'AS15169', router: 'core-rtr-01', status: 'running', time: 'running' },
  { id: 'q4', type: 'ping', target: '1.1.1.1', router: 'edge-rtr-03', status: 'completed', time: '12 min ago' },
];

const routerHealth = [
  { name: 'core-rtr-01', vendor: 'MikroTik', status: 'online', load: 45, uptime: '99.98%' },
  { name: 'edge-rtr-02', vendor: 'Juniper', status: 'online', load: 62, uptime: '99.95%' },
  { name: 'edge-rtr-03', vendor: 'Cisco', status: 'online', load: 38, uptime: '99.99%' },
  { name: 'core-rtr-04', vendor: 'Huawei', status: 'warning', load: 78, uptime: '99.87%' },
];

export default function DashboardPage() {
  const { data: statsData, isLoading } = useQuery({
    queryKey: ['admin-stats'],
    queryFn: async () => {
      const response = await adminApi.stats();
      return response.data;
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">
          Welcome back! Here&apos;s an overview of your network diagnostics platform.
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <Card key={stat.name}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">{stat.name}</CardTitle>
                <Icon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                <p className="text-xs text-muted-foreground flex items-center">
                  {stat.trend === 'up' ? (
                    <ArrowUpRight className="mr-1 h-4 w-4 text-green-500" />
                  ) : (
                    <ArrowDownRight className="mr-1 h-4 w-4 text-red-500" />
                  )}
                  {stat.change} from last month
                </p>
              </CardContent>
            </Card>
          );
        })}
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Recent Queries</CardTitle>
            <CardDescription>Latest network diagnostic queries</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentQueries.map((query) => (
                <div
                  key={query.id}
                  className="flex items-center justify-between border-b pb-3 last:border-0"
                >
                  <div className="flex items-center gap-3">
                    <div className="rounded-lg bg-muted p-2">
                      <Globe className="h-4 w-4" />
                    </div>
                    <div>
                      <p className="font-medium">{query.target}</p>
                      <p className="text-sm text-muted-foreground">
                        {query.type} • {query.router}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge
                      variant={
                        query.status === 'completed'
                          ? 'success'
                          : query.status === 'running'
                          ? 'warning'
                          : 'destructive'
                      }
                    >
                      {query.status}
                    </Badge>
                    <span className="text-sm text-muted-foreground">{query.time}</span>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Router Health</CardTitle>
            <CardDescription>Current status of your routers</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {routerHealth.map((router) => (
                <div key={router.name} className="space-y-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <div
                        className={`h-2 w-2 rounded-full ${
                          router.status === 'online'
                            ? 'bg-green-500'
                            : router.status === 'warning'
                            ? 'bg-yellow-500'
                            : 'bg-red-500'
                        }`}
                      />
                      <span className="font-medium">{router.name}</span>
                      <Badge variant="outline" className="text-xs">
                        {router.vendor}
                      </Badge>
                    </div>
                    <span className="text-sm text-muted-foreground">
                      {router.uptime}
                    </span>
                  </div>
                  <Progress value={router.load} className="h-2" />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-2">
            <Link href="/looking-glass">
              <Button className="w-full justify-start" variant="outline">
                <Globe className="mr-2 h-4 w-4" />
                Run Network Query
              </Button>
            </Link>
            <Link href="/routers">
              <Button className="w-full justify-start" variant="outline">
                <Server className="mr-2 h-4 w-4" />
                Manage Routers
              </Button>
            </Link>
            <Link href="/admin">
              <Button className="w-full justify-start" variant="outline">
                <Users className="mr-2 h-4 w-4" />
                View All Users
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Query Types Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm">Ping</span>
                <span className="font-medium">45%</span>
              </div>
              <Progress value={45} />
              <div className="flex items-center justify-between">
                <span className="text-sm">Traceroute</span>
                <span className="font-medium">30%</span>
              </div>
              <Progress value={30} />
              <div className="flex items-center justify-between">
                <span className="text-sm">BGP Lookup</span>
                <span className="font-medium">25%</span>
              </div>
              <Progress value={25} />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>System Status</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm">Database</span>
              <Badge variant="success">Connected</Badge>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm">Redis Cache</span>
              <Badge variant="success">Connected</Badge>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm">SSH Pool</span>
              <Badge variant="success">{24}/{24} Active</Badge>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm">Queue Workers</span>
              <Badge variant="success">5 Running</Badge>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}