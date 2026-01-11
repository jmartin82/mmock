import { useState } from 'react';
import { cn } from '@/lib/utils';
import { ChevronDown, ChevronRight, Copy, Check } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface JsonViewerProps {
  data: string | object;
  className?: string;
  maxHeight?: string;
}

export const JsonViewer = ({ data, className, maxHeight = '400px' }: JsonViewerProps) => {
  const [copied, setCopied] = useState(false);

  let parsedData: unknown;
  let isValidJson = true;
  let rawString = '';

  if (typeof data === 'string') {
    rawString = data;
    try {
      parsedData = JSON.parse(data);
    } catch {
      isValidJson = false;
      parsedData = data;
    }
  } else {
    parsedData = data;
    rawString = JSON.stringify(data, null, 2);
  }

  const handleCopy = async () => {
    await navigator.clipboard.writeText(rawString);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  if (!isValidJson) {
    return (
      <div className={cn('relative', className)}>
        <Button
          variant="ghost"
          size="icon"
          className="absolute right-2 top-2 h-7 w-7"
          onClick={handleCopy}
        >
          {copied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
        </Button>
        <pre
          className="bg-muted/50 rounded-md p-4 font-mono text-sm overflow-auto scrollbar-thin"
          style={{ maxHeight }}
        >
          {rawString || <span className="text-muted-foreground italic">Empty</span>}
        </pre>
      </div>
    );
  }

  return (
    <div className={cn('relative', className)}>
      <Button
        variant="ghost"
        size="icon"
        className="absolute right-2 top-2 h-7 w-7 z-10"
        onClick={handleCopy}
      >
        {copied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
      </Button>
      <div
        className="bg-muted/50 rounded-md p-4 font-mono text-sm overflow-auto scrollbar-thin"
        style={{ maxHeight }}
      >
        <JsonNode data={parsedData} />
      </div>
    </div>
  );
};

interface JsonNodeProps {
  data: unknown;
  depth?: number;
}

const JsonNode = ({ data, depth = 0 }: JsonNodeProps) => {
  const [isExpanded, setIsExpanded] = useState(depth < 2);

  if (data === null) {
    return <span className="json-null">null</span>;
  }

  if (typeof data === 'boolean') {
    return <span className="json-boolean">{data.toString()}</span>;
  }

  if (typeof data === 'number') {
    return <span className="json-number">{data}</span>;
  }

  if (typeof data === 'string') {
    return <span className="json-string">"{data}"</span>;
  }

  if (Array.isArray(data)) {
    if (data.length === 0) {
      return <span className="json-bracket">[]</span>;
    }

    return (
      <span>
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="inline-flex items-center hover:bg-accent/50 rounded"
        >
          {isExpanded ? (
            <ChevronDown className="h-3 w-3 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-3 w-3 text-muted-foreground" />
          )}
        </button>
        <span className="json-bracket">[</span>
        {!isExpanded && (
          <span className="text-muted-foreground">
            {data.length} item{data.length !== 1 ? 's' : ''}
          </span>
        )}
        {isExpanded && (
          <div className="ml-4">
            {data.map((item, index) => (
              <div key={index}>
                <JsonNode data={item} depth={depth + 1} />
                {index < data.length - 1 && <span className="json-bracket">,</span>}
              </div>
            ))}
          </div>
        )}
        <span className="json-bracket">]</span>
      </span>
    );
  }

  if (typeof data === 'object') {
    const entries = Object.entries(data);
    if (entries.length === 0) {
      return <span className="json-bracket">{'{}'}</span>;
    }

    return (
      <span>
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="inline-flex items-center hover:bg-accent/50 rounded"
        >
          {isExpanded ? (
            <ChevronDown className="h-3 w-3 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-3 w-3 text-muted-foreground" />
          )}
        </button>
        <span className="json-bracket">{'{'}</span>
        {!isExpanded && (
          <span className="text-muted-foreground">
            {entries.length} key{entries.length !== 1 ? 's' : ''}
          </span>
        )}
        {isExpanded && (
          <div className="ml-4">
            {entries.map(([key, value], index) => (
              <div key={key}>
                <span className="json-key">"{key}"</span>
                <span className="json-bracket">: </span>
                <JsonNode data={value} depth={depth + 1} />
                {index < entries.length - 1 && <span className="json-bracket">,</span>}
              </div>
            ))}
          </div>
        )}
        <span className="json-bracket">{'}'}</span>
      </span>
    );
  }

  return null;
};
