import { useState, useEffect } from 'react';
import { RefreshCw, RotateCcw, Play, Pause, GitBranch, Edit2, Check } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Input } from '@/components/ui/input';
import { scenarioApi, Scenarios as ScenariosType } from '@/lib/api';
import { toast } from 'sonner';

export default function Scenarios() {
  const [scenarios, setScenarios] = useState<ScenariosType>({});
  const [loading, setLoading] = useState(true);
  const [isPaused, setIsPaused] = useState(false);
  const [editingScenario, setEditingScenario] = useState<string | null>(null);
  const [newState, setNewState] = useState('');

  const fetchScenarios = async () => {
    setLoading(true);
    try {
      const data = await scenarioApi.getAll();
      setScenarios(data || {});
    } catch (error) {
      toast.error('Failed to fetch scenarios');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchScenarios();
  }, []);

  const handleResetAll = async () => {
    try {
      await scenarioApi.resetAll();
      toast.success('All scenarios reset');
      fetchScenarios();
    } catch (error) {
      toast.error('Failed to reset scenarios');
      console.error(error);
    }
  };

  const handleTogglePause = async () => {
    try {
      if (isPaused) {
        await scenarioApi.unpause();
        toast.success('Scenarios resumed');
      } else {
        await scenarioApi.pause();
        toast.success('Scenarios paused');
      }
      setIsPaused(!isPaused);
    } catch (error) {
      toast.error('Failed to toggle pause state');
      console.error(error);
    }
  };

  const handleSetState = async (scenario: string) => {
    if (!newState.trim()) {
      toast.error('Please enter a state name');
      return;
    }
    try {
      await scenarioApi.setState(scenario, newState);
      toast.success(`Scenario "${scenario}" set to "${newState}"`);
      setEditingScenario(null);
      setNewState('');
      fetchScenarios();
    } catch (error) {
      toast.error('Failed to set scenario state');
      console.error(error);
    }
  };

  const scenarioEntries = Object.entries(scenarios);

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="border-b border-border bg-card/30 p-4">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <h1 className="text-xl font-semibold">Scenarios</h1>
            <Badge variant="secondary" className="font-mono">
              {scenarioEntries.length} scenarios
            </Badge>
            {isPaused && (
              <Badge variant="outline" className="text-status-4xx border-status-4xx/30 bg-status-4xx/10">
                <Pause className="h-3 w-3 mr-1" />
                Paused
              </Badge>
            )}
          </div>

          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={fetchScenarios} disabled={loading}>
              <RefreshCw className={`h-4 w-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
            <Button variant="outline" size="sm" onClick={handleTogglePause}>
              {isPaused ? (
                <>
                  <Play className="h-4 w-4 mr-2" />
                  Resume
                </>
              ) : (
                <>
                  <Pause className="h-4 w-4 mr-2" />
                  Pause
                </>
              )}
            </Button>
            <Button variant="outline" size="sm" onClick={handleResetAll}>
              <RotateCcw className="h-4 w-4 mr-2" />
              Reset All
            </Button>
          </div>
        </div>
      </div>

      {/* Scenarios List */}
      <div className="flex-1 overflow-auto scrollbar-thin p-4">
        {loading ? (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <Card key={i}>
                <CardHeader>
                  <Skeleton className="h-5 w-2/3" />
                </CardHeader>
                <CardContent>
                  <Skeleton className="h-8 w-full" />
                </CardContent>
              </Card>
            ))}
          </div>
        ) : scenarioEntries.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center p-8">
            <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4">
              <GitBranch className="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 className="text-lg font-medium mb-2">No scenarios configured</h3>
            <p className="text-muted-foreground max-w-md">
              Scenarios allow you to manage complex test flows with state machines. Configure scenarios in your mock definitions to see them here.
            </p>
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {scenarioEntries.map(([name, state]) => (
              <Card key={name} className="hover:border-primary/50 transition-colors">
                <CardHeader className="pb-3">
                  <CardTitle className="text-base font-mono flex items-center gap-2">
                    <GitBranch className="h-4 w-4 text-primary" />
                    {name}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {editingScenario === name ? (
                    <div className="flex items-center gap-2">
                      <Input
                        placeholder="New state..."
                        value={newState}
                        onChange={(e) => setNewState(e.target.value)}
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') handleSetState(name);
                          if (e.key === 'Escape') {
                            setEditingScenario(null);
                            setNewState('');
                          }
                        }}
                        className="flex-1"
                        autoFocus
                      />
                      <Button size="icon" onClick={() => handleSetState(name)}>
                        <Check className="h-4 w-4" />
                      </Button>
                    </div>
                  ) : (
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="text-sm text-muted-foreground">Current State:</span>
                        <Badge variant="secondary" className="font-mono">
                          {state || 'not_started'}
                        </Badge>
                      </div>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={() => {
                          setEditingScenario(name);
                          setNewState(state);
                        }}
                      >
                        <Edit2 className="h-4 w-4" />
                      </Button>
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
