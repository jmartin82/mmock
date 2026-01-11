import { useState, useEffect } from 'react';
import { Save, Server, Check } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { getBaseUrl, setBaseUrl } from '@/lib/api';
import { toast } from 'sonner';

export default function Settings() {
  const [apiUrl, setApiUrl] = useState('');
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    setApiUrl(getBaseUrl());
  }, []);

  const handleSave = () => {
    if (!apiUrl.trim()) {
      toast.error('API URL cannot be empty');
      return;
    }

    try {
      new URL(apiUrl);
      setBaseUrl(apiUrl);
      toast.success('Settings saved. Reload the page to connect to the new server.');
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch {
      toast.error('Please enter a valid URL');
    }
  };

  const handleReset = () => {
    const defaultUrl = window.location.protocol + "//" + window.location.host;
    setApiUrl(defaultUrl);
    setBaseUrl(defaultUrl);
    toast.success('Reset to default URL');
  };

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold">Settings</h1>
        <p className="text-muted-foreground mt-1">Configure your MMock Console connection</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Server className="h-5 w-5" />
            API Configuration
          </CardTitle>
          <CardDescription>
            Configure the MMock server URL. This is used for all API calls and WebSocket connections.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="apiUrl">MMock Server URL</Label>
            <Input
              id="apiUrl"
              placeholder={window.location.protocol + "//" + window.location.host}
              value={apiUrl}
              onChange={(e) => setApiUrl(e.target.value)}
              className="font-mono"
            />
            <p className="text-xs text-muted-foreground">
              The base URL of your MMock server. WebSocket connections will use the same host.
            </p>
          </div>

          <div className="flex items-center gap-3 pt-2">
            <Button onClick={handleSave} disabled={saved} className="gap-2">
              {saved ? (
                <>
                  <Check className="h-4 w-4" />
                  Saved
                </>
              ) : (
                <>
                  <Save className="h-4 w-4" />
                  Save Settings
                </>
              )}
            </Button>
            <Button variant="outline" onClick={handleReset}>
              Reset to Default
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Connection Tips</CardTitle>
        </CardHeader>
        <CardContent className="text-sm text-muted-foreground space-y-3">
          <p>
            <strong>HTTPS/WSS:</strong> If your MMock server uses HTTPS, the console will automatically
            use secure WebSocket (WSS) for real-time updates.
          </p>
          <p>
            <strong>CORS:</strong> Ensure your MMock server allows CORS requests from this console's
            origin if running on different domains.
          </p>
          <p>
            <strong>Self-signed Certificates:</strong> If using self-signed certificates for HTTPS,
            you may need to visit the MMock server URL directly and accept the certificate first.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
