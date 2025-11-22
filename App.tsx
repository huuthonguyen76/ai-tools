import React, { useState, useCallback } from 'react';
import { Sparkles, Link as LinkIcon, ExternalLink, AlertCircle, Quote, FileText, Copy, Check, Hash, ArrowRight } from 'lucide-react';
import { Background } from './components/Background';
import { ResultCard } from './components/ResultCard';
import { generateLinkContext } from './services/geminiService';
import { GenerationState } from './types';

const App: React.FC = () => {
  const [url, setUrl] = useState('');
  const [state, setState] = useState<GenerationState>({
    isLoading: false,
    error: null,
    result: null,
  });

  const handleAnalyze = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url.trim()) return;

    let validUrl = url.trim();
    if (!validUrl.startsWith('http')) {
        validUrl = `https://${validUrl}`;
    }

    setState(prev => ({ ...prev, isLoading: true, error: null, result: null }));

    try {
      const data = await generateLinkContext(validUrl);
      setState({ isLoading: false, error: null, result: data });
    } catch (err) {
      setState({ 
        isLoading: false, 
        error: err instanceof Error ? err.message : "Failed to contextualize link.", 
        result: null 
      });
    }
  }, [url]);

  // Generates the Clean Readable Link (Slug style): example.com/page#my-cool-summary
  const getReadableLink = () => {
    if (!state.result) return '';
    const { originalUrl, descriptiveLabel } = state.result;
    
    // Create a slug: Lowercase, remove special chars, replace spaces with dashes
    const slug = descriptiveLabel
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, '') // Remove non-alphanumeric chars except spaces and dashes
      .trim()
      .replace(/\s+/g, '-');        // Replace spaces with dashes
      
    return `${originalUrl}#${slug}`;
  };

  // Generates the Technical Deep Link (Text Fragment): example.com/page#:~:text=...
  const getDeepLink = () => {
    if (!state.result) return '';
    const { originalUrl, highlightPhrase } = state.result;
    if (!highlightPhrase) return originalUrl;
    const encodedPhrase = encodeURIComponent(highlightPhrase);
    return `${originalUrl}#:~:text=${encodedPhrase}`;
  };

  return (
    <div className="min-h-screen relative font-sans text-slate-200">
      <Background />
      
      <main className="container mx-auto px-4 py-12 md:py-20 max-w-4xl">
        
        {/* Header */}
        <div className="text-center mb-12 space-y-4">
          <div className="inline-flex items-center justify-center p-1.5 mb-2 rounded-full bg-indigo-500/10 border border-indigo-500/20 backdrop-blur-sm">
            <span className="px-3 py-1 text-xs font-bold text-indigo-300 tracking-wide uppercase">
              Contextual Link Engine
            </span>
          </div>
          <h1 className="text-4xl md:text-5xl font-bold tracking-tight text-white">
            Contextualize Your <span className="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-cyan-400">Link</span>
          </h1>
          <p className="text-lg text-slate-400 max-w-xl mx-auto">
            Transform raw URLs into descriptive, readable smart links. No more messy codes—just clear, English context.
          </p>
        </div>

        {/* Input Section */}
        <div className="max-w-xl mx-auto mb-16 relative z-10">
          <form onSubmit={handleAnalyze} className="relative group">
            <div className="absolute -inset-1 bg-gradient-to-r from-emerald-500 to-cyan-500 rounded-2xl blur opacity-20 group-hover:opacity-40 transition duration-500"></div>
            <div className="relative flex items-center bg-slate-900 rounded-xl border border-slate-700 shadow-2xl overflow-hidden">
              <div className="pl-5 text-slate-500">
                <LinkIcon className="w-5 h-5" />
              </div>
              <input
                type="text"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="Paste any website link here..."
                className="w-full bg-transparent border-none py-4 px-4 text-base text-white placeholder-slate-500 focus:outline-none focus:ring-0"
              />
              <button
                type="submit"
                disabled={state.isLoading || !url}
                className="mr-1.5 px-5 py-2.5 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-medium transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 text-sm"
              >
                {state.isLoading ? (
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                ) : (
                  <>
                    <span>Contextualize</span>
                    <Sparkles className="w-4 h-4" />
                  </>
                )}
              </button>
            </div>
          </form>
          
          {state.error && (
            <div className="mt-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg flex items-center gap-3 text-red-400 text-sm animate-in fade-in slide-in-from-top-2">
              <AlertCircle className="w-4 h-4 shrink-0" />
              <p>{state.error}</p>
            </div>
          )}
        </div>

        {/* Results */}
        {state.result && (
          <div className="animate-in fade-in duration-700 space-y-8">
            
            {/* Main Result Box */}
            <div className="bg-gradient-to-b from-slate-800/80 to-slate-900/80 border border-slate-700/50 rounded-2xl p-6 md:p-8 backdrop-blur-md shadow-2xl ring-1 ring-white/5">
               
               {/* 1. Context Summary */}
               <div className="border-b border-slate-700/50 pb-6 mb-6">
                  <div className="flex items-center gap-2 mb-3">
                    <div className="h-6 w-1 bg-emerald-500 rounded-full"></div>
                    <h2 className="text-sm font-semibold text-emerald-400 uppercase tracking-wider">About This Link</h2>
                  </div>
                  <p className="text-xl md:text-2xl font-medium text-white leading-relaxed">
                    {state.result.summary}
                  </p>
                  <div className="mt-4 flex flex-wrap gap-2">
                    {state.result.suggestedTags.map((tag, i) => (
                      <span key={i} className="px-2.5 py-1 rounded-md bg-slate-700/50 border border-slate-600/50 text-xs text-slate-300">
                        #{tag.toLowerCase().replace(/\s+/g, '-')}
                      </span>
                    ))}
                  </div>
               </div>

               {/* 2. The Generated Links */}
               <div className="space-y-6">
                  
                  {/* Primary: Readable Smart Link */}
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-indigo-300 mb-1">
                        <Hash className="w-4 h-4" />
                        <span className="text-xs font-bold uppercase tracking-wider">Readable Smart Link</span>
                    </div>
                    <ResultCard 
                        label="Contextual URL" 
                        content={getReadableLink()} 
                        isCode 
                        delay={100}
                    />
                    <p className="text-xs text-slate-500 px-1">
                        A clean link containing the context in the URL itself. Great for sharing in emails and messages.
                    </p>
                  </div>

                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 pt-4 border-t border-slate-800/50">
                    {/* Markdown & HTML using Readable Link */}
                    <div className="space-y-4">
                        <h3 className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Copy Code</h3>
                        <ResultCard 
                        label="Markdown" 
                        content={`[${state.result.descriptiveLabel}](${getReadableLink()})`} 
                        isCode 
                        delay={200}
                        />
                        <ResultCard 
                        label="HTML" 
                        content={`<a href="${getReadableLink()}">${state.result.descriptiveLabel}</a>`} 
                        isCode 
                        delay={300}
                        />
                    </div>

                     {/* Secondary: Deep Link */}
                    <div className="space-y-4">
                        <h3 className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Precision Deep Link</h3>
                        <ResultCard 
                            label="Highlight URL" 
                            content={getDeepLink()} 
                            isCode 
                            delay={400}
                        />
                        <div className="p-3 rounded-lg bg-slate-800/50 border border-slate-700/50 text-xs text-slate-400">
                            <div className="flex items-start gap-2">
                                <Quote className="w-3 h-3 mt-0.5 shrink-0 text-emerald-500" />
                                <span>
                                    highlights: <span className="italic text-emerald-400/80">"{state.result.highlightPhrase}"</span>
                                </span>
                            </div>
                        </div>
                    </div>
                  </div>

               </div>

               <div className="mt-8 flex items-center justify-between text-xs text-slate-500 pt-4 border-t border-slate-800">
                 <div className="flex gap-2">
                    <span>Original Source: {state.result.title}</span>
                 </div>
                 <a 
                    href={getReadableLink()} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className="text-indigo-400 hover:text-indigo-300 flex items-center gap-1 transition-colors font-medium group"
                 >
                   Test Smart Link <ArrowRight className="w-3 h-3 group-hover:translate-x-0.5 transition-transform" />
                 </a>
               </div>
            </div>

          </div>
        )}
      </main>
    </div>
  );
};

export default App;