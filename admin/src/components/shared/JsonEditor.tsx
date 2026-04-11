import { Alert, AlertDescription, AlertIcon, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { WandSparkles } from 'lucide-react';

export type JSONEditorProps = {
  id: string;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
};

export function getJSONError(value: string) {
  try {
    JSON.parse(value);
    return '';
  } catch (error) {
    const message = error instanceof Error ? error.message : 'JSON 格式无效';
    const matched = message.match(/position\s+(\d+)/i);
    if (!matched) {
      return message;
    }

    const position = Number(matched[1]);
    if (Number.isNaN(position)) {
      return message;
    }

    const cursor = value.slice(0, position);
    const line = cursor.split('\n').length;
    const column = cursor.length - cursor.lastIndexOf('\n');
    return `第 ${line} 行，第 ${column} 列附近格式有误：${message}`;
  }
}

export function formatJSONValue(value: string) {
  return JSON.stringify(JSON.parse(value), null, 2);
}

export function JSONEditor({ id, value, onChange, disabled = false }: JSONEditorProps) {
  const error = getJSONError(value);

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between gap-3">
        <div className="text-sm text-muted-foreground">支持实时校验和一键格式化。</div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={disabled || !!error}
          onClick={() => onChange(formatJSONValue(value))}
        >
          <WandSparkles className="h-4 w-4" />
          格式化 JSON
        </Button>
      </div>

      <Textarea
        id={id}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        className="min-h-[260px] font-mono text-xs leading-6"
        spellCheck={false}
        disabled={disabled}
        required
      />

      {error ? (
        <Alert variant="destructive">
          <AlertIcon className="h-4 w-4" />
          <AlertTitle>JSON 格式无效</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      ) : (
        <div className="text-sm text-muted-foreground">当前 JSON 格式有效，可以直接保存。</div>
      )}
    </div>
  );
}
