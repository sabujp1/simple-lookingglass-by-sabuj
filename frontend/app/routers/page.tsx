'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { 
  Server, 
  Plus, 
  Search, 
  MoreHorizontal, 
  Pencil, 
  Trash2, 
  RefreshCw,
  CheckCircle,
  XCircle,
  Loader2
} from 'lucide-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { routerApi, zoneApi } from '@/lib/api';

interface Router {
  id: string;
  name: string;
  hostname: string;
  ip_address: string;
  vendor: string;
  model: string;
  port: number;
  username: string;
  is_active: boolean;
  zone_id: string;
  zone_name?: string;
  last_check?: string;
  created_at: string;
}

interface Zone {
  id: string;
  name: string;
}

export default function RoutersPage() {
  const [search, setSearch] = useState('');
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [editingRouter, setEditingRouter] = useState<Router | null>(null);

  const queryClient = useQueryClient();

  const { data: routersData, isLoading } = useQuery({
    queryKey: ['routers'],
    queryFn: async () => {
      const response = await routerApi.list({ limit: 100 });
      return response.data;
    },
  });

  const { data: zonesData } = useQuery({
    queryKey: ['zones'],
    queryFn: async () => {
      const response = await zoneApi.list({ limit: 100 });
      return response.data;
    },
  });

  const zones = zonesData?.data || [];
  const routers = routersData?.data || [];

  const filteredRouters = routers.filter((router: Router) =>
    router.name.toLowerCase().includes(search.toLowerCase()) ||
    router.hostname.toLowerCase().includes(search.toLowerCase()) ||
    router.ip_address.includes(search)
  );

  const testMutation = useMutation({
    mutationFn: async (id: string) => {
      const response = await routerApi.test(id);
      return response.data;
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await routerApi.delete(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['routers'] });
    },
  });

  const toggleActiveMutation = useMutation({
    mutationFn: async ({ id, is_active }: { id: string; is_active: boolean }) => {
      const response = await routerApi.update(id, { is_active });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['routers'] });
    },
  });

  const getVendorColor = (vendor: string) => {
    switch (vendor.toLowerCase()) {
      case 'mikrotik':
        return 'bg-orange-100 text-orange-800';
      case 'juniper':
        return 'bg-red-100 text-red-800';
      case 'cisco':
        return 'bg-blue-100 text-blue-800';
      case 'huawei':
        return 'bg-red-50 text-red-700';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Routers</h1>
          <p className="text-muted-foreground">
            Manage your network routers and their connections
          </p>
        </div>
        <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Add Router
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Add New Router</DialogTitle>
              <DialogDescription>
                Configure a new router for network diagnostics
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Name</Label>
                  <Input id="name" placeholder="core-rtr-01" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="zone">Zone</Label>
                  <Select>
                    <SelectTrigger>
                      <SelectValue placeholder="Select zone" />
                    </SelectTrigger>
                    <SelectContent>
                      {zones.map((zone: Zone) => (
                        <SelectItem key={zone.id} value={zone.id}>
                          {zone.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="hostname">Hostname</Label>
                  <Input id="hostname" placeholder="router1.example.com" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="ip">IP Address</Label>
                  <Input id="ip" placeholder="192.168.1.1" />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="vendor">Vendor</Label>
                  <Select>
                    <SelectTrigger>
                      <SelectValue placeholder="Select vendor" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="mikrotik">MikroTik</SelectItem>
                      <SelectItem value="juniper">Juniper</SelectItem>
                      <SelectItem value="cisco">Cisco</SelectItem>
                      <SelectItem value="huawei">Huawei</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="model">Model</Label>
                  <Input id="model" placeholder="CCR-1072" />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="username">Username</Label>
                  <Input id="username" placeholder="admin" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="port">SSH Port</Label>
                  <Input id="port" type="number" placeholder="22" defaultValue={22} />
                </div>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsAddOpen(false)}>
                Cancel
              </Button>
              <Button onClick={() => setIsAddOpen(false)}>
                Add Router
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>All Routers</CardTitle>
              <CardDescription>
                {filteredRouters.length} router{filteredRouters.length !== 1 ? 's' : ''} configured
              </CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <div className="relative">
                <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search routers..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-8 w-64"
                />
              </div>
              <Button variant="outline" size="icon">
                <RefreshCw className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>IP Address</TableHead>
                <TableHead>Vendor</TableHead>
                <TableHead>Zone</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Active</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8">
                    <Loader2 className="h-6 w-6 animate-spin mx-auto" />
                  </TableCell>
                </TableRow>
              ) : filteredRouters.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                    No routers found
                  </TableCell>
                </TableRow>
              ) : (
                filteredRouters.map((router: Router) => (
                  <TableRow key={router.id}>
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        <Server className="h-4 w-4 text-muted-foreground" />
                        {router.name}
                      </div>
                    </TableCell>
                    <TableCell className="font-mono text-sm">{router.ip_address}</TableCell>
                    <TableCell>
                      <Badge className={getVendorColor(router.vendor)}>
                        {router.vendor}
                      </Badge>
                    </TableCell>
                    <TableCell>{router.zone_name || '-'}</TableCell>
                    <TableCell>
                      {router.is_active ? (
                        <div className="flex items-center gap-1 text-green-600">
                          <CheckCircle className="h-4 w-4" />
                          <span className="text-sm">Online</span>
                        </div>
                      ) : (
                        <div className="flex items-center gap-1 text-red-600">
                          <XCircle className="h-4 w-4" />
                          <span className="text-sm">Offline</span>
                        </div>
                      )}
                    </TableCell>
                    <TableCell>
                      <Switch
                        checked={router.is_active}
                        onCheckedChange={(checked) =>
                          toggleActiveMutation.mutate({ id: router.id, is_active: checked })
                        }
                      />
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => testMutation.mutate(router.id)}
                          disabled={testMutation.isPending}
                        >
                          <RefreshCw className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => setEditingRouter(router)}
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="text-destructive hover:text-destructive"
                          onClick={() => deleteMutation.mutate(router.id)}
                          disabled={deleteMutation.isPending}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}