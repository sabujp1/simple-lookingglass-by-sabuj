'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Globe, Server, Activity, TrendingUp, Play, Loader2, CheckCircle2, XCircle, Clock } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { routerApi, queryApi } from '@/lib/api';
import { useWebSocket } from '@/lib/websocket';

interface QueryResult {
  id: string;
  status: string;
  output: string;
  error?: string;
}

export default function LookingGlassPage() {
  const [queryType, setQueryType] = useState<'ping' | 'traceroute' | 'bgp'>('ping');
  const [target, setTarget] = useState('');
  const [selectedRouter, setSelectedRouter] = useState('');
  const [queryResult, setQueryResult] = useState<QueryResult | null>(null);
  const [isExecuting, setIsExecuting] = useState(false);

  const { data: routersData } = useQuery({
    queryKey: ['routers'],
    queryFn: async () => {
      const response = await routerApi.list({ limit: 100 });
      return response.data;
    },
  });

  const routers = routersData?.data || [];

  const { connected, messages, send } = useWebSocket(
    queryResult?.id || null,
    {
      onMessage: (msg) => {
        if (msg.type === 'output' && msg.data) {
          setQueryResult((prev) => ({
            ...prev!,
            output: prev?.output + msg.data + '\n',
          }));
        } else if (msg.type === 'complete') {
          setQueryResult((prev) => ({ ...prev!, status: 'completed' }));
        } else if (msg.type === 'error') {
          setQueryResult((prev) => ({
            ...prev!,
            status: 'failed',
            error: msg.error,
          }));
        }
      },
    }
  );

  const executeQuery = async () => {
    if (!selectedRouter || !target) return;

    setIsExecuting(true);
    setQueryResult({ id: '', status: 'pending', output: '' });

    try {
      let response;
      switch (queryType) {
        case 'ping':
          response = await queryApi.ping(selectedRouter, target);
          break;
        case 'traceroute':
          response = await queryApi.traceroute(selectedRouter, target);
          break;
        case 'bgp':
          response = await queryApi.bgp(selectedRouter, target);
          break;
      }
      setQueryResult({
        ...response.data,
        status: response.data.status === 'pending' ? 'running' : response.data.status,
      });
    } catch (error) {
      setQueryResult({
        id: '',
        status: 'failed',
        output: '',
        error: 'Failed to execute query',
      });
    } finally {
      setIsExecuting(false);
    }
  };

  const getRouterVendorIcon = (vendor: string) => {
    switch (vendor.toLowerCase()) {
      case 'mikrotik':
        return <Server className="h-4 w-4 text-orange-500" />;
      case 'juniper':
        return <Server className="h-4 w-4 text-red-500" />;
      case 'cisco':
        return <Server className="h-4 w-4 text-blue-500" />;
      case 'huawei':
        return <Server className="h-4 w-4 text-red-600" />;
      default:
        return <Server className="h-4 w-4" />;
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Looking Glass</h1>
        <p className="text-muted-foreground">
          Execute network diagnostic queries on your routers
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Query Interface</CardTitle>
            <CardDescription>
              Select a router and query type to execute network diagnostics
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="router">Router</Label>
                <Select value={selectedRouter} onValueChange={setSelectedRouter}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select a router" />
                  </SelectTrigger>
                  <SelectContent>
                    {routers.map((router: any) => (
                      <SelectItem key={router.id} value={router.id}>
                        <div className="flex items-center gap-2">
                          {getRouterVendorIcon(router.vendor)}
                          <span>{router.name}</span>
                          <Badge variant="outline" className="text-xs">
                            {router.vendor}
                          </Badge>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="query-type">Query Type</Label>
                <Select value={queryType} onValueChange={(v: any) => setQueryType(v)}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="ping">
                      <div className="flex items-center gap-2">
                        <Activity className="h-4 w-4" />
                        Ping
                      </div>
                    </SelectItem>
                    <SelectItem value="traceroute">
                      <div className="flex items-center gap-2">
                        <TrendingUp className="h-4 w-4" />
                        Traceroute
                      </div>
                    </SelectItem>
                    <SelectItem value="bgp">
                      <div className="flex items-center gap-2">
                        <Globe className="h-4 w-4" />
                        BGP Lookup
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="target">
                {queryType === 'ping' && 'Target IP/Hostname'}
                {queryType === 'traceroute' && 'Target IP/Hostname'}
                {queryType === 'bgp' && 'IP Address or AS Number'}
              </Label>
              <Input
                id="target"
                placeholder={
                  queryType === 'ping' ? '8.8.8.8 or google.com' :
                  queryType === 'traceroute' ? '1.1.1.1 or cloudflare.com' :
                  '8.8.8.8 or AS15169'
                }
                value={target}
                onChange={(e) => setTarget(e.target.value)}
              />
            </div>

            <Button
              onClick={executeQuery}
              disabled={!selectedRouter || !target || isExecuting}
              className="w-full"
            >
              {isExecuting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Executing...
                </>
              ) : (
                <>
                  <Play className="mr-2 h-4 w-4" />
                  Execute Query
                </>
              )}
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Query Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="rounded-lg border p-4">
              <h4 className="font-medium mb-2">Query Types</h4>
              <div className="space-y-2 text-sm">
                <div className="flex items-center gap-2">
                  <Activity className="h-4 w-4 text-muted-foreground" />
                  <span><strong>Ping:</strong> Test connectivity and latency</span>
                </div>
                <div className="flex items-center gap-2">
                  <TrendingUp className="h-4 w-4 text-muted-foreground" />
                  <span><strong>Traceroute:</strong> Trace network path</span>
                </div>
                <div className="flex items-center gap-2">
                  <Globe className="h-4 w-4 text-muted-foreground" />
                  <span><strong>BGP:</strong> View routing information</span>
                </div>
              </div>
            </div>

            <div className="rounded-lg border p-4">
              <h4 className="font-medium mb-2">Supported Vendors</h4>
              <div className="flex flex-wrap gap-2">
                <Badge variant="outline">MikroTik</Badge>
                <Badge variant="outline">Juniper</Badge>
                <Badge variant="outline">Cisco</Badge>
                <Badge variant="outline">Huawei</Badge>
              </div>
            </div>

            {queryResult && (
              <div className="rounded-lg border p-4">
                <h4 className="font-medium mb-2">Status</h4>
                <div className="flex items-center gap-2">
                  {queryResult.status === 'completed' && (
                    <>
                      <CheckCircle2 className="h-4 w-4 text-green-500" />
                      <span className="text-sm text-green-500">Completed</span>
                    </>
                  )}
                  {queryResult.status === 'failed' && (
                    <>
                      <XCircle className="h-4 w-4 text-red-500" />
                      <span className="text-sm text-red-500">Failed</span>
                    </>
                  )}
                  {queryResult.status === 'running' && (
                    <>
                      <Clock className="h-4 w-4 text-yellow-500 animate-pulse" />
                      <span className="text-sm text-yellow-500">Running...</span>
                    </>
                  )}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {queryResult && (
        <Card>
          <CardHeader>
            <CardTitle>Query Result</CardTitle>
            <CardDescription>
              Output from {selectedRouter ? routers.find((r: any) => r.id === selectedRouter)?.name : 'router'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <pre className="font-mono text-sm bg-muted p-4 rounded-lg overflow-x-auto whitespace-pre-wrap">
              {queryResult.output || queryResult.error || 'Waiting for output...'}
            </pre>
          </CardContent>
        </Card>
      )}
    </div>
  );
}