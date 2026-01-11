import { useState, useMemo, useEffect } from 'react';
import { Trash2, RefreshCw, Filter, Search, Database } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { TransactionRow } from '@/components/TransactionRow';
import { requestApi, Transaction } from '@/lib/api';
import { toast } from 'sonner';
import { Skeleton } from '@/components/ui/skeleton';

type FilterType = 'all' | 'matched' | 'unmatched';

export default function RequestMemory() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<FilterType>('all');
  const [searchQuery, setSearchQuery] = useState('');

  const fetchTransactions = async () => {
    setLoading(true);
    try {
      const data = await requestApi.getAll();
      setTransactions(data || []);
    } catch (error) {
      toast.error('Failed to fetch requests');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTransactions();
  }, []);

  const filteredTransactions = useMemo(() => {
    let filtered = transactions;

    if (filter === 'matched') {
      filtered = filtered.filter((t) => t.result.match);
    } else if (filter === 'unmatched') {
      filtered = filtered.filter((t) => !t.result.match);
    }

    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (t) =>
          t.request.path.toLowerCase().includes(query) ||
          t.request.method.toLowerCase().includes(query) ||
          t.response.statusCode.toString().includes(query)
      );
    }

    return filtered;
  }, [transactions, filter, searchQuery]);

  const handleClearRequests = async () => {
    try {
      await requestApi.reset();
      setTransactions([]);
      toast.success('All requests cleared');
    } catch (error) {
      toast.error('Failed to clear requests');
      console.error(error);
    }
  };

  const matchedCount = transactions.filter((t) => t.result.match).length;
  const unmatchedCount = transactions.filter((t) => !t.result.match).length;

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="border-b border-border bg-card/30 p-4">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <h1 className="text-xl font-semibold">Request Memory</h1>
            <Badge variant="secondary" className="gap-1.5">
              <Database className="h-3 w-3" />
              Stored Requests
            </Badge>
          </div>

          <div className="flex items-center gap-2">
            <Badge variant="secondary" className="font-mono">
              {transactions.length} total
            </Badge>
            <Badge variant="outline" className="font-mono text-status-2xx border-status-2xx/30">
              {matchedCount} matched
            </Badge>
            <Badge variant="outline" className="font-mono text-status-4xx border-status-4xx/30">
              {unmatchedCount} unmatched
            </Badge>
          </div>
        </div>

        {/* Controls */}
        <div className="flex flex-col md:flex-row md:items-center gap-4 mt-4">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search by path, method, or status..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9 bg-background"
            />
          </div>

          <Tabs value={filter} onValueChange={(v) => setFilter(v as FilterType)}>
            <TabsList>
              <TabsTrigger value="all" className="gap-2">
                <Filter className="h-3.5 w-3.5" />
                All
              </TabsTrigger>
              <TabsTrigger value="matched">Matched</TabsTrigger>
              <TabsTrigger value="unmatched">Unmatched</TabsTrigger>
            </TabsList>
          </Tabs>

          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={fetchTransactions} disabled={loading}>
              <RefreshCw className={`h-4 w-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
            <Button variant="outline" size="sm" onClick={handleClearRequests}>
              <Trash2 className="h-4 w-4 mr-2" />
              Clear
            </Button>
          </div>
        </div>
      </div>

      {/* Request List */}
      <div className="flex-1 overflow-auto scrollbar-thin">
        {loading ? (
          <div className="p-4 space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="flex items-center gap-4 p-3">
                <Skeleton className="h-4 w-4" />
                <Skeleton className="h-6 w-16" />
                <Skeleton className="h-4 flex-1" />
                <Skeleton className="h-4 w-12" />
                <Skeleton className="h-6 w-20" />
                <Skeleton className="h-4 w-16" />
              </div>
            ))}
          </div>
        ) : filteredTransactions.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center p-8">
            <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4">
              <Database className="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 className="text-lg font-medium mb-2">No stored requests</h3>
            <p className="text-muted-foreground max-w-md">
              {searchQuery
                ? 'No requests match your search criteria. Try adjusting your filters.'
                : 'No requests have been stored yet. Make requests to your mock server to see them here.'}
            </p>
          </div>
        ) : (
          <div className="divide-y divide-border">
            {filteredTransactions.map((transaction, index) => (
              <TransactionRow
                key={`${transaction.time}-${index}`}
                transaction={transaction}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
