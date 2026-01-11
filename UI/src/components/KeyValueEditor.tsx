import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Plus, Trash2 } from 'lucide-react';

interface KeyValueEditorProps {
  label: string;
  value: Record<string, string> | Record<string, string[]> | undefined;
  onChange: (value: Record<string, string> | Record<string, string[]> | undefined) => void;
  arrayValues?: boolean; // If true, values are string arrays; if false, single strings
  placeholder?: { key?: string; value?: string };
}

export function KeyValueEditor({
  label,
  value,
  onChange,
  arrayValues = false,
  placeholder = { key: 'Key', value: 'Value' },
}: KeyValueEditorProps) {
  const entries = Object.entries(value || {});

  const addEntry = () => {
    const newKey = `key_${Date.now()}`;
    if (arrayValues) {
      onChange({ ...(value as Record<string, string[]>), [newKey]: [''] });
    } else {
      onChange({ ...(value as Record<string, string>), [newKey]: '' });
    }
  };

  const removeEntry = (keyToRemove: string) => {
    const newValue = { ...value };
    delete newValue[keyToRemove];
    onChange(Object.keys(newValue).length > 0 ? newValue : undefined);
  };

  const updateKey = (oldKey: string, newKey: string) => {
    if (oldKey === newKey) return;
    const newValue: Record<string, string | string[]> = {};
    for (const [k, v] of Object.entries(value || {})) {
      if (k === oldKey) {
        newValue[newKey] = v;
      } else {
        newValue[k] = v;
      }
    }
    onChange(newValue as Record<string, string> | Record<string, string[]>);
  };

  const updateValue = (key: string, newVal: string | string[]) => {
    onChange({ ...value, [key]: newVal } as Record<string, string> | Record<string, string[]>);
  };

  const formatArrayValue = (val: string | string[]): string => {
    if (Array.isArray(val)) {
      return val.join(', ');
    }
    return val;
  };

  const parseArrayValue = (val: string): string[] => {
    return val.split(',').map((s) => s.trim()).filter((s) => s.length > 0);
  };

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <Label className="text-sm font-medium">{label}</Label>
        <Button type="button" variant="outline" size="sm" onClick={addEntry} className="h-7 gap-1">
          <Plus className="h-3 w-3" />
          Add
        </Button>
      </div>

      {entries.length === 0 ? (
        <p className="text-xs text-muted-foreground italic">No {label.toLowerCase()} configured</p>
      ) : (
        <div className="space-y-2">
          {entries.map(([key, val], index) => (
            <div key={index} className="flex gap-2 items-start">
              <Input
                placeholder={placeholder.key}
                value={key.startsWith('key_') ? '' : key}
                onChange={(e) => updateKey(key, e.target.value || `key_${Date.now()}_${index}`)}
                className="flex-1 h-8 text-sm"
              />
              <Input
                placeholder={arrayValues ? `${placeholder.value} (comma-separated)` : placeholder.value}
                value={formatArrayValue(val)}
                onChange={(e) =>
                  arrayValues
                    ? updateValue(key, parseArrayValue(e.target.value) || [''])
                    : updateValue(key, e.target.value)
                }
                className="flex-[2] h-8 text-sm"
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => removeEntry(key)}
                className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
              >
                <Trash2 className="h-3.5 w-3.5" />
              </Button>
            </div>
          ))}
        </div>
      )}
      {arrayValues && entries.length > 0 && (
        <p className="text-xs text-muted-foreground">Multiple values can be separated by commas</p>
      )}
    </div>
  );
}
