import { useState, useEffect } from "react";
import { Plus, Search, RefreshCw, Trash2, Edit, FileJson } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { MethodBadge } from "@/components/MethodBadge";
import { StatusBadge } from "@/components/StatusBadge";
import { mappingApi, MockDefinition } from "@/lib/api";
import { toast } from "sonner";
import { MockEditor } from "@/components/MockEditor";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

export default function MockMappings() {
  const [mappings, setMappings] = useState<MockDefinition[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [editorOpen, setEditorOpen] = useState(false);
  const [editingMock, setEditingMock] = useState<MockDefinition | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<string | null>(null);

  const fetchMappings = async () => {
    setLoading(true);
    try {
      const data = await mappingApi.getAll();
      setMappings(data || []);
    } catch (error) {
      toast.error("Failed to fetch mock mappings");
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMappings();
  }, []);

  const filteredMappings = mappings.filter((mock) => {
    if (!searchQuery) return true;
    const query = searchQuery.toLowerCase();
    return (
      mock.URI.toLowerCase().includes(query) ||
      mock.description?.toLowerCase().includes(query) ||
      mock.request.path.toLowerCase().includes(query) ||
      mock.request.method.toLowerCase().includes(query)
    );
  });

  const handleCreate = () => {
    setEditingMock(null);
    setEditorOpen(true);
  };

  const handleEdit = (mock: MockDefinition) => {
    setEditingMock(mock);
    setEditorOpen(true);
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    try {
      await mappingApi.delete(deleteTarget);
      toast.success("Mock deleted successfully");
      fetchMappings();
    } catch (error) {
      toast.error("Failed to delete mock");
      console.error(error);
    } finally {
      setDeleteTarget(null);
    }
  };

  const handleSave = async (mock: MockDefinition) => {
    try {
      if (editingMock) {
        await mappingApi.update(mock.URI, mock);
        toast.success("Mock updated successfully");
      } else {
        await mappingApi.create(mock.URI, mock);
        toast.success("Mock created successfully");
      }
      setEditorOpen(false);
      fetchMappings();
    } catch (error) {
      toast.error("Failed to save mock");
      console.error(error);
    }
  };

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="border-b border-border bg-card/30 p-4">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <h1 className="text-xl font-semibold">Mock Mappings</h1>
            <Badge variant="secondary" className="font-mono">
              {mappings.length} definitions
            </Badge>
          </div>

          <Button onClick={handleCreate} className="gap-2">
            <Plus className="h-4 w-4" />
            Add Mock
          </Button>
        </div>

        {/* Controls */}
        <div className="flex flex-col md:flex-row md:items-center gap-4 mt-4">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search mocks..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9 bg-background"
            />
          </div>

          <Button variant="outline" size="sm" onClick={fetchMappings} disabled={loading}>
            <RefreshCw className={`h-4 w-4 mr-2 ${loading ? "animate-spin" : ""}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Mapping List */}
      <div className="flex-1 overflow-auto scrollbar-thin p-4">
        {loading ? (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <Card key={i}>
                <CardContent className="p-4">
                  <div className="space-y-3">
                    <Skeleton className="h-5 w-3/4" />
                    <Skeleton className="h-4 w-full" />
                    <div className="flex gap-2">
                      <Skeleton className="h-6 w-16" />
                      <Skeleton className="h-6 w-24" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : filteredMappings.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center p-8">
            <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4">
              <FileJson className="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 className="text-lg font-medium mb-2">No mock mappings</h3>
            <p className="text-muted-foreground max-w-md mb-4">
              {searchQuery
                ? "No mocks match your search criteria."
                : "Create your first mock definition to start mocking HTTP endpoints."}
            </p>
            {!searchQuery && (
              <Button onClick={handleCreate} className="gap-2">
                <Plus className="h-4 w-4" />
                Add Mock
              </Button>
            )}
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {filteredMappings.map((mock) => (
              <Card key={mock.URI} className="hover:border-primary/50 transition-colors">
                <CardContent className="p-4">
                  <div className="flex items-start justify-between gap-2 mb-3">
                    <div className="min-w-0 flex-1">
                      <h3 className="font-mono text-sm font-medium truncate" title={mock.URI}>
                        {mock.URI}
                      </h3>
                      {mock.description && (
                        <p className="text-sm text-muted-foreground truncate mt-1">{mock.description}</p>
                      )}
                    </div>
                    <div className="flex items-center gap-1 flex-shrink-0">
                      <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => handleEdit(mock)}>
                        <Edit className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-destructive hover:text-destructive"
                        onClick={() => setDeleteTarget(mock.URI)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>

                  <div className="flex items-center gap-2 flex-wrap">
                    <MethodBadge method={mock.request.method} />
                    <code className="text-xs bg-muted px-2 py-1 rounded truncate max-w-[200px]">
                      {mock.request.path}
                    </code>
                  </div>

                  <div className="flex items-center gap-2 mt-3 pt-3 border-t border-border">
                    <span className="text-xs text-muted-foreground">Response:</span>
                    <StatusBadge statusCode={mock.response.statusCode} />
                    {mock.control?.delay && (
                      <Badge variant="outline" className="text-xs">
                        Delay: {mock.control.delay}
                      </Badge>
                    )}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>

      {/* Mock Editor Dialog */}
      <MockEditor open={editorOpen} onOpenChange={setEditorOpen} mock={editingMock} onSave={handleSave} />

      {/* Delete Confirmation */}
      <AlertDialog open={!!deleteTarget} onOpenChange={() => setDeleteTarget(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Mock Definition</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this mock definition? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
