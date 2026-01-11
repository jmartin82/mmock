import { useState } from 'react';
import { ChevronDown, ChevronRight, Check, X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Transaction } from '@/lib/api';
import { MethodBadge } from './MethodBadge';
import { StatusBadge } from './StatusBadge';
import { JsonViewer } from './JsonViewer';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';

interface TransactionRowProps {
  transaction: Transaction;
  className?: string;
}

export const TransactionRow = ({ transaction, className }: TransactionRowProps) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const { request, response, result, time } = transaction;
  const matched = result.match;

  const formatTime = (ts: number) => {
    const date = new Date(ts * 1000);
    return date.toLocaleTimeString('en-US', {
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  const hasQueryParams = Object.keys(request.queryStringParameters || {}).length > 0;
  const hasRequestHeaders = Object.keys(request.headers || {}).length > 0;
  const hasRequestCookies = Object.keys(request.cookies || {}).length > 0;
  const hasResponseHeaders = Object.keys(response.headers || {}).length > 0;
  const hasResponseCookies = Object.keys(response.cookies || {}).length > 0;

  return (
    <div
      className={cn(
        'border-b border-border hover:bg-accent/30 transition-colors animate-fade-in',
        className
      )}
    >
      <button
        className="w-full flex items-center gap-4 px-4 py-3 text-left"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <div className="flex-shrink-0">
          {isExpanded ? (
            <ChevronDown className="h-4 w-4 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-4 w-4 text-muted-foreground" />
          )}
        </div>

        <MethodBadge method={request.method} className="w-16 flex-shrink-0" />

        <div className="flex-1 min-w-0">
          <span className="font-mono text-sm truncate block">{request.path}</span>
        </div>

        <StatusBadge statusCode={response.statusCode} className="flex-shrink-0" />

        <div className="flex-shrink-0">
          {matched ? (
            <Badge variant="outline" className="text-status-2xx border-status-2xx/30 bg-status-2xx/10">
              <Check className="h-3 w-3 mr-1" />
              Matched
            </Badge>
          ) : (
            <Badge variant="outline" className="text-status-4xx border-status-4xx/30 bg-status-4xx/10">
              <X className="h-3 w-3 mr-1" />
              Unmatched
            </Badge>
          )}
        </div>

        <span className="text-xs text-muted-foreground font-mono flex-shrink-0 w-20 text-right">
          {formatTime(time)}
        </span>
      </button>

      {isExpanded && (
        <div className="px-4 pb-4 animate-slide-in">
          <Tabs defaultValue="request" className="w-full">
            <TabsList className="mb-4">
              <TabsTrigger value="request">Request</TabsTrigger>
              <TabsTrigger value="response">Response</TabsTrigger>
              <TabsTrigger value="match">Match Info</TabsTrigger>
            </TabsList>

            <TabsContent value="request" className="space-y-4">
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                <div>
                  <span className="text-muted-foreground">Scheme</span>
                  <p className="font-mono">{request.scheme}</p>
                </div>
                <div>
                  <span className="text-muted-foreground">Host</span>
                  <p className="font-mono">{request.host}</p>
                </div>
                <div>
                  <span className="text-muted-foreground">Port</span>
                  <p className="font-mono">{request.port}</p>
                </div>
                <div>
                  <span className="text-muted-foreground">Method</span>
                  <MethodBadge method={request.method} />
                </div>
              </div>

              <div>
                <span className="text-muted-foreground text-sm">Path</span>
                <p className="font-mono text-sm bg-muted/50 p-2 rounded-md mt-1">
                  {request.path}
                </p>
              </div>

              {hasQueryParams && (
                <div>
                  <span className="text-muted-foreground text-sm">Query Parameters</span>
                  <JsonViewer data={request.queryStringParameters} maxHeight="200px" />
                </div>
              )}

              {hasRequestHeaders && (
                <div>
                  <span className="text-muted-foreground text-sm">Headers</span>
                  <JsonViewer data={request.headers} maxHeight="200px" />
                </div>
              )}

              {hasRequestCookies && (
                <div>
                  <span className="text-muted-foreground text-sm">Cookies</span>
                  <JsonViewer data={request.cookies} maxHeight="200px" />
                </div>
              )}

              {request.body && (
                <div>
                  <span className="text-muted-foreground text-sm">Body</span>
                  <JsonViewer data={request.body} maxHeight="300px" />
                </div>
              )}
            </TabsContent>

            <TabsContent value="response" className="space-y-4">
              <div className="flex items-center gap-4">
                <span className="text-muted-foreground text-sm">Status Code</span>
                <StatusBadge statusCode={response.statusCode} />
              </div>

              {hasResponseHeaders && (
                <div>
                  <span className="text-muted-foreground text-sm">Headers</span>
                  <JsonViewer data={response.headers} maxHeight="200px" />
                </div>
              )}

              {hasResponseCookies && (
                <div>
                  <span className="text-muted-foreground text-sm">Cookies</span>
                  <JsonViewer data={response.cookies} maxHeight="200px" />
                </div>
              )}

              {response.body && (
                <div>
                  <span className="text-muted-foreground text-sm">Body</span>
                  <JsonViewer data={response.body} maxHeight="300px" />
                </div>
              )}
            </TabsContent>

            <TabsContent value="match" className="space-y-4">
              <div className="flex items-center gap-4">
                <span className="text-muted-foreground text-sm">Match Status</span>
                {matched ? (
                  <Badge className="bg-status-2xx/20 text-status-2xx hover:bg-status-2xx/30">
                    <Check className="h-3 w-3 mr-1" />
                    Matched
                  </Badge>
                ) : (
                  <Badge className="bg-status-4xx/20 text-status-4xx hover:bg-status-4xx/30">
                    <X className="h-3 w-3 mr-1" />
                    No Match
                  </Badge>
                )}
              </div>

              {matched && result.uri && (
                <div>
                  <span className="text-muted-foreground text-sm">Matched Mock</span>
                  <p className="font-mono text-sm bg-status-2xx/10 p-2 rounded-md mt-1 text-status-2xx">
                    {result.uri}
                  </p>
                </div>
              )}

              {result.errors && result.errors.length > 0 && (
                <div>
                  <span className="text-muted-foreground text-sm">
                    Match Attempts ({result.errors.length} mocks checked)
                  </span>
                  <div className="mt-2 max-h-[300px] overflow-y-auto space-y-2">
                    {result.errors.map((error, index) => (
                      <div 
                        key={index} 
                        className="flex items-start gap-3 p-2 bg-muted/50 rounded-md text-sm"
                      >
                        <X className="h-4 w-4 text-status-4xx flex-shrink-0 mt-0.5" />
                        <div className="min-w-0">
                          <span className="font-mono font-medium block truncate">{error.uri}</span>
                          <p className="text-muted-foreground">{error.reason}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {!matched && (!result.errors || result.errors.length === 0) && (
                <div className="text-sm text-muted-foreground">
                  No matching mock found. Consider adding a mock for this endpoint.
                </div>
              )}
            </TabsContent>
          </Tabs>
        </div>
      )}
    </div>
  );
};
