import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { JSONEditor } from './JsonEditor';

describe('JSONEditor', () => {
  it('非法 JSON 时展示错误提示', async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();

    render(<JSONEditor id="json-editor" value="{invalid" onChange={onChange} />);

    expect(screen.getByText('JSON 格式无效')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '格式化 JSON' })).toBeDisabled();

    await user.type(screen.getByRole('textbox'), '}');
    expect(onChange).toHaveBeenCalled();
  });

  it('格式化按钮会输出格式化后的 JSON', async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();

    render(<JSONEditor id="json-editor" value='{"name":"gin-template"}' onChange={onChange} />);

    await user.click(screen.getByRole('button', { name: '格式化 JSON' }));

    expect(onChange).toHaveBeenCalledWith('{\n  "name": "gin-template"\n}');
  });
});
