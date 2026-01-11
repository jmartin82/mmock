import { useState, useEffect } from 'react';
import { MockDefinition } from '@/lib/api';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Switch } from '@/components/ui/switch';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import { Save, ChevronDown } from 'lucide-react';
import { KeyValueEditor } from '@/components/KeyValueEditor';

interface MockEditorProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  mock: MockDefinition | null;
  onSave: (mock: MockDefinition) => void;
}

const defaultMock: MockDefinition = {
  URI: '',
  description: '',
  request: {
    method: 'GET',
    path: '',
  },
  response: {
    statusCode: 200,
  },
};

const HTTP_METHODS = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'];

export function MockEditor({ open, onOpenChange, mock, onSave }: MockEditorProps) {
  const [formData, setFormData] = useState<MockDefinition>(defaultMock);
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (mock) {
      setFormData(mock);
    } else {
      setFormData(defaultMock);
    }
    setErrors({});
  }, [mock, open]);

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.URI.trim()) {
      newErrors.URI = 'URI is required';
    }
    if (!formData.request.path.trim()) {
      newErrors.path = 'Request path is required';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = () => {
    if (validate()) {
      onSave(formData);
    }
  };

  const updateRequest = (field: string, value: string | Record<string, string> | Record<string, string[]> | undefined) => {
    setFormData((prev) => ({
      ...prev,
      request: { ...prev.request, [field]: value },
    }));
  };

  const updateResponse = (field: string, value: string | number | Record<string, string> | Record<string, string[]> | undefined) => {
    setFormData((prev) => ({
      ...prev,
      response: { ...prev.response, [field]: value },
    }));
  };

  const updateControl = (field: string, value: string | number | boolean | undefined) => {
    setFormData((prev) => ({
      ...prev,
      control: { ...prev.control, [field]: value },
    }));
  };

  const updateScenario = (field: string, value: string | string[] | undefined) => {
    setFormData((prev) => ({
      ...prev,
      control: {
        ...prev.control,
        scenario: {
          ...prev.control?.scenario,
          name: prev.control?.scenario?.name || '',
          [field]: value,
        },
      },
    }));
  };

  const updateCallback = (field: string, value: string | Record<string, string> | Record<string, string[]> | undefined) => {
    setFormData((prev) => ({
      ...prev,
      callback: { ...prev.callback, [field]: value },
    }));
  };

  const parseRequiredStates = (value: string): string[] => {
    return value.split(',').map(s => s.trim()).filter(s => s.length > 0);
  };

  const formatRequiredStates = (states?: string[]): string => {
    return states?.join(', ') || '';
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] flex flex-col overflow-hidden p-0">
        <DialogHeader className="flex-shrink-0 p-6 pb-0">
          <DialogTitle>{mock ? 'Edit Mock Definition' : 'Create Mock Definition'}</DialogTitle>
        </DialogHeader>

        <div className="flex-1 overflow-hidden px-6">
          <ScrollArea className="h-[calc(90vh-180px)]">
            <div className="space-y-6 pb-4 pr-4">
            {/* Basic Info */}
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="uri">URI *</Label>
                  <Input
                    id="uri"
                    placeholder="unique-mock-id"
                    value={formData.URI}
                    onChange={(e) => setFormData((prev) => ({ ...prev, URI: e.target.value }))}
                    className={errors.URI ? 'border-destructive' : ''}
                    disabled={!!mock}
                  />
                  {errors.URI && <p className="text-xs text-destructive">{errors.URI}</p>}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Input
                    id="description"
                    placeholder="Mock description"
                    value={formData.description}
                    onChange={(e) => setFormData((prev) => ({ ...prev, description: e.target.value }))}
                  />
                </div>
              </div>
            </div>

            {/* Request Section */}
            <Accordion type="single" collapsible defaultValue="request">
              <AccordionItem value="request">
                <AccordionTrigger className="text-sm font-medium">Request</AccordionTrigger>
                <AccordionContent className="space-y-4 pt-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label>Method *</Label>
                      <Select
                        value={formData.request.method}
                        onValueChange={(v) => updateRequest('method', v)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {HTTP_METHODS.map((method) => (
                            <SelectItem key={method} value={method}>
                              {method}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="path">Path *</Label>
                      <Input
                        id="path"
                        placeholder="/api/users/:id"
                        value={formData.request.path}
                        onChange={(e) => updateRequest('path', e.target.value)}
                        className={errors.path ? 'border-destructive' : ''}
                      />
                      {errors.path && <p className="text-xs text-destructive">{errors.path}</p>}
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="scheme">Scheme</Label>
                      <Input
                        id="scheme"
                        placeholder="http"
                        value={formData.request.scheme || ''}
                        onChange={(e) => updateRequest('scheme', e.target.value || undefined)}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="host">Host</Label>
                      <Input
                        id="host"
                        placeholder="example.com"
                        value={formData.request.host || ''}
                        onChange={(e) => updateRequest('host', e.target.value || undefined)}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="port">Port</Label>
                      <Input
                        id="port"
                        placeholder="8080"
                        value={formData.request.port || ''}
                        onChange={(e) => updateRequest('port', e.target.value || undefined)}
                      />
                    </div>
                  </div>

                  {/* Query String Parameters */}
                  <Collapsible className="rounded-lg border border-border">
                    <CollapsibleTrigger className="flex w-full items-center justify-between p-3 hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Query String Parameters</span>
                        {formData.request.queryStringParameters && Object.keys(formData.request.queryStringParameters).length > 0 && (
                          <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded">
                            {Object.keys(formData.request.queryStringParameters).length}
                          </span>
                        )}
                      </div>
                      <ChevronDown className="h-4 w-4 text-muted-foreground" />
                    </CollapsibleTrigger>
                    <CollapsibleContent className="border-t border-border p-3">
                      <KeyValueEditor
                        label="Parameters"
                        value={formData.request.queryStringParameters}
                        onChange={(v) => updateRequest('queryStringParameters', v)}
                        arrayValues={true}
                        placeholder={{ key: 'name', value: 'value1, value2' }}
                      />
                    </CollapsibleContent>
                  </Collapsible>

                  {/* Request Headers */}
                  <Collapsible className="rounded-lg border border-border">
                    <CollapsibleTrigger className="flex w-full items-center justify-between p-3 hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Headers</span>
                        {formData.request.headers && Object.keys(formData.request.headers).length > 0 && (
                          <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded">
                            {Object.keys(formData.request.headers).length}
                          </span>
                        )}
                      </div>
                      <ChevronDown className="h-4 w-4 text-muted-foreground" />
                    </CollapsibleTrigger>
                    <CollapsibleContent className="border-t border-border p-3">
                      <KeyValueEditor
                        label="Headers"
                        value={formData.request.headers}
                        onChange={(v) => updateRequest('headers', v)}
                        arrayValues={true}
                        placeholder={{ key: 'Content-Type', value: 'application/json' }}
                      />
                    </CollapsibleContent>
                  </Collapsible>

                  {/* Request Cookies */}
                  <Collapsible className="rounded-lg border border-border">
                    <CollapsibleTrigger className="flex w-full items-center justify-between p-3 hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Cookies</span>
                        {formData.request.cookies && Object.keys(formData.request.cookies).length > 0 && (
                          <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded">
                            {Object.keys(formData.request.cookies).length}
                          </span>
                        )}
                      </div>
                      <ChevronDown className="h-4 w-4 text-muted-foreground" />
                    </CollapsibleTrigger>
                    <CollapsibleContent className="border-t border-border p-3">
                      <KeyValueEditor
                        label="Cookies"
                        value={formData.request.cookies}
                        onChange={(v) => updateRequest('cookies', v as Record<string, string> | undefined)}
                        arrayValues={false}
                        placeholder={{ key: 'session_id', value: 'abc123' }}
                      />
                    </CollapsibleContent>
                  </Collapsible>

                  <div className="space-y-2">
                    <Label htmlFor="requestBody">Request Body</Label>
                    <Textarea
                      id="requestBody"
                      placeholder="JSON, text, XML, etc."
                      value={formData.request.body || ''}
                      onChange={(e) => updateRequest('body', e.target.value || undefined)}
                      className="font-mono text-sm min-h-[100px]"
                    />
                  </div>
                </AccordionContent>
              </AccordionItem>

              <AccordionItem value="response">
                <AccordionTrigger className="text-sm font-medium">Response</AccordionTrigger>
                <AccordionContent className="space-y-4 pt-4">
                  <div className="space-y-2">
                    <Label htmlFor="statusCode">Status Code</Label>
                    <Input
                      id="statusCode"
                      type="number"
                      placeholder="200"
                      value={formData.response.statusCode}
                      onChange={(e) => updateResponse('statusCode', parseInt(e.target.value) || 200)}
                      className="w-32"
                    />
                  </div>

                  {/* Response Headers */}
                  <Collapsible className="rounded-lg border border-border">
                    <CollapsibleTrigger className="flex w-full items-center justify-between p-3 hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Headers</span>
                        {formData.response.headers && Object.keys(formData.response.headers).length > 0 && (
                          <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded">
                            {Object.keys(formData.response.headers).length}
                          </span>
                        )}
                      </div>
                      <ChevronDown className="h-4 w-4 text-muted-foreground" />
                    </CollapsibleTrigger>
                    <CollapsibleContent className="border-t border-border p-3">
                      <KeyValueEditor
                        label="Headers"
                        value={formData.response.headers}
                        onChange={(v) => updateResponse('headers', v)}
                        arrayValues={true}
                        placeholder={{ key: 'Content-Type', value: 'application/json' }}
                      />
                    </CollapsibleContent>
                  </Collapsible>

                  {/* Response Cookies */}
                  <Collapsible className="rounded-lg border border-border">
                    <CollapsibleTrigger className="flex w-full items-center justify-between p-3 hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Cookies</span>
                        {formData.response.cookies && Object.keys(formData.response.cookies).length > 0 && (
                          <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded">
                            {Object.keys(formData.response.cookies).length}
                          </span>
                        )}
                      </div>
                      <ChevronDown className="h-4 w-4 text-muted-foreground" />
                    </CollapsibleTrigger>
                    <CollapsibleContent className="border-t border-border p-3">
                      <KeyValueEditor
                        label="Cookies"
                        value={formData.response.cookies}
                        onChange={(v) => updateResponse('cookies', v as Record<string, string> | undefined)}
                        arrayValues={false}
                        placeholder={{ key: 'session_id', value: 'abc123' }}
                      />
                    </CollapsibleContent>
                  </Collapsible>

                  <div className="space-y-2">
                    <Label htmlFor="responseBody">Response Body</Label>
                    <Textarea
                      id="responseBody"
                      placeholder="JSON, text, XML, etc."
                      value={formData.response.body || ''}
                      onChange={(e) => updateResponse('body', e.target.value || undefined)}
                      className="font-mono text-sm min-h-[150px]"
                    />
                  </div>
                </AccordionContent>
              </AccordionItem>

              <AccordionItem value="control">
                <AccordionTrigger className="text-sm font-medium">Control (Advanced)</AccordionTrigger>
                <AccordionContent className="space-y-4 pt-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="priority">Priority</Label>
                      <Input
                        id="priority"
                        type="number"
                        placeholder="0"
                        value={formData.control?.priority || ''}
                        onChange={(e) => updateControl('priority', parseInt(e.target.value) || undefined)}
                      />
                      <p className="text-xs text-muted-foreground">Higher priority mocks are matched first</p>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="delay">Delay</Label>
                      <Input
                        id="delay"
                        placeholder="100ms or 1s"
                        value={formData.control?.delay || ''}
                        onChange={(e) => updateControl('delay', e.target.value || undefined)}
                      />
                      <p className="text-xs text-muted-foreground">Add latency to response</p>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="proxyUrl">Proxy Base URL</Label>
                      <Input
                        id="proxyUrl"
                        placeholder="https://api.example.com"
                        value={formData.control?.proxyBaseURL || ''}
                        onChange={(e) => updateControl('proxyBaseURL', e.target.value || undefined)}
                      />
                      <p className="text-xs text-muted-foreground">Forward requests to real server</p>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="webHookUrl">Webhook URL</Label>
                      <Input
                        id="webHookUrl"
                        placeholder="https://webhook.example.com"
                        value={formData.control?.webHookURL || ''}
                        onChange={(e) => updateControl('webHookURL', e.target.value || undefined)}
                      />
                      <p className="text-xs text-muted-foreground">Notify external service</p>
                    </div>
                  </div>

                  {/* Crazy Mode Toggle */}
                  <div className="flex items-center justify-between rounded-lg border border-border p-4">
                    <div className="space-y-0.5">
                      <Label htmlFor="crazy" className="text-base">Crazy Mode</Label>
                      <p className="text-xs text-muted-foreground">Return random 5xx server errors</p>
                    </div>
                    <Switch
                      id="crazy"
                      checked={formData.control?.crazy || false}
                      onCheckedChange={(checked) => updateControl('crazy', checked)}
                    />
                  </div>

                  {/* Scenario Section */}
                  <Collapsible className="rounded-lg border border-border">
                    <CollapsibleTrigger className="flex w-full items-center justify-between p-4 hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Scenario Configuration</span>
                        {formData.control?.scenario?.name && (
                          <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded">
                            {formData.control.scenario.name}
                          </span>
                        )}
                      </div>
                      <ChevronDown className="h-4 w-4 text-muted-foreground transition-transform duration-200 group-data-[state=open]:rotate-180" />
                    </CollapsibleTrigger>
                    <CollapsibleContent className="border-t border-border p-4 space-y-4">
                      <div className="space-y-2">
                        <Label htmlFor="scenarioName">Scenario Name</Label>
                        <Input
                          id="scenarioName"
                          placeholder="my-scenario"
                          value={formData.control?.scenario?.name || ''}
                          onChange={(e) => updateScenario('name', e.target.value || undefined)}
                        />
                        <p className="text-xs text-muted-foreground">Unique identifier for the scenario state machine</p>
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="requiredState">Required States</Label>
                        <Input
                          id="requiredState"
                          placeholder="not_started, state_1, state_2"
                          value={formatRequiredStates(formData.control?.scenario?.requiredState)}
                          onChange={(e) => updateScenario('requiredState', parseRequiredStates(e.target.value))}
                        />
                        <p className="text-xs text-muted-foreground">Comma-separated list of states required to match this mock</p>
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="newState">New State</Label>
                        <Input
                          id="newState"
                          placeholder="next_state"
                          value={formData.control?.scenario?.newState || ''}
                          onChange={(e) => updateScenario('newState', e.target.value || undefined)}
                        />
                        <p className="text-xs text-muted-foreground">State to transition to after this mock is matched</p>
                      </div>
                    </CollapsibleContent>
                  </Collapsible>
                </AccordionContent>
              </AccordionItem>

              {/* Callback Section */}
              <AccordionItem value="callback">
                <AccordionTrigger className="text-sm font-medium">Callback (Advanced)</AccordionTrigger>
                <AccordionContent className="space-y-4 pt-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label>Method</Label>
                      <Select
                        value={formData.callback?.method || ''}
                        onValueChange={(v) => updateCallback('method', v || undefined)}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Select method" />
                        </SelectTrigger>
                        <SelectContent>
                          {HTTP_METHODS.map((method) => (
                            <SelectItem key={method} value={method}>
                              {method}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="callbackUrl">URL</Label>
                      <Input
                        id="callbackUrl"
                        placeholder="http://your-callback/"
                        value={formData.callback?.url || ''}
                        onChange={(e) => updateCallback('url', e.target.value || undefined)}
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="callbackDelay">Delay</Label>
                    <Input
                      id="callbackDelay"
                      placeholder="100ms or 1s"
                      value={formData.callback?.delay || ''}
                      onChange={(e) => updateCallback('delay', e.target.value || undefined)}
                      className="w-48"
                    />
                    <p className="text-xs text-muted-foreground">Delay before sending callback</p>
                  </div>

                  {/* Callback Headers */}
                  <Collapsible className="rounded-lg border border-border">
                    <CollapsibleTrigger className="flex w-full items-center justify-between p-3 hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Headers</span>
                        {formData.callback?.headers && Object.keys(formData.callback.headers).length > 0 && (
                          <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded">
                            {Object.keys(formData.callback.headers).length}
                          </span>
                        )}
                      </div>
                      <ChevronDown className="h-4 w-4 text-muted-foreground" />
                    </CollapsibleTrigger>
                    <CollapsibleContent className="border-t border-border p-3">
                      <KeyValueEditor
                        label="Headers"
                        value={formData.callback?.headers}
                        onChange={(v) => updateCallback('headers', v)}
                        arrayValues={true}
                        placeholder={{ key: 'Content-Type', value: 'application/json' }}
                      />
                    </CollapsibleContent>
                  </Collapsible>

                  {/* Callback Cookies */}
                  <Collapsible className="rounded-lg border border-border">
                    <CollapsibleTrigger className="flex w-full items-center justify-between p-3 hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Cookies</span>
                        {formData.callback?.cookies && Object.keys(formData.callback.cookies).length > 0 && (
                          <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded">
                            {Object.keys(formData.callback.cookies).length}
                          </span>
                        )}
                      </div>
                      <ChevronDown className="h-4 w-4 text-muted-foreground" />
                    </CollapsibleTrigger>
                    <CollapsibleContent className="border-t border-border p-3">
                      <KeyValueEditor
                        label="Cookies"
                        value={formData.callback?.cookies}
                        onChange={(v) => updateCallback('cookies', v as Record<string, string> | undefined)}
                        arrayValues={false}
                        placeholder={{ key: 'session_id', value: 'abc123' }}
                      />
                    </CollapsibleContent>
                  </Collapsible>

                  <div className="space-y-2">
                    <Label htmlFor="callbackBody">Body</Label>
                    <Textarea
                      id="callbackBody"
                      placeholder="JSON, text, XML, etc."
                      value={formData.callback?.body || ''}
                      onChange={(e) => updateCallback('body', e.target.value || undefined)}
                      className="font-mono text-sm min-h-[100px]"
                    />
                  </div>
                </AccordionContent>
              </AccordionItem>
            </Accordion>
          </div>
          </ScrollArea>
        </div>

        <div className="flex justify-end gap-3 p-6 pt-4 border-t border-border flex-shrink-0">
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleSubmit} className="gap-2">
            <Save className="h-4 w-4" />
            {mock ? 'Update Mock' : 'Create Mock'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
