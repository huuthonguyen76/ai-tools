import React, { useState } from 'react';
import { Copy, Check } from 'lucide-react';

interface ResultCardProps {
  label: string;
  content: string;
  isCode?: boolean;
  delay?: number;
}

export const ResultCard: React.FC<ResultCardProps> = ({ label, content, isCode = false, delay = 0 }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(content);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div 
      className="group relative bg-slate-900/60 border border-slate-700/50 rounded-xl p-4 hover:bg-slate-800/80 transition-all duration-300 animate-in fade-in slide-in-from-bottom-4 fill-mode-backwards"
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="flex items-center justify-between mb-2">
        <label className="text-xs font-semibold text-indigo-300 uppercase tracking-wider">{label}</label>
        <button
          onClick={handleCopy}
          className="flex items-center gap-1.5 px-2 py-1 text-xs font-medium text-slate-400 bg-slate-800 hover:bg-indigo-600 hover:text-white rounded-md transition-colors"
        >
          {copied ? (
            <>
              <Check className="w-3 h-3" />
              <span>Copied</span>
            </>
          ) : (
            <>
              <Copy className="w-3 h-3" />
              <span>Copy</span>
            </>
          )}
        </button>
      </div>
      
      <div className={`w-full break-all ${isCode ? 'font-mono text-xs text-emerald-300/90' : 'text-sm text-slate-200'}`}>
        {content}
      </div>
    </div>
  );
};