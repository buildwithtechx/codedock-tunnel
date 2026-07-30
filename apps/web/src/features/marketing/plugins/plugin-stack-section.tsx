import { motion } from 'motion/react';
import { MarketingContainer } from '#/components/layout';
import type { PluginDefinition } from './plugin-data';

export function PluginStackSection({ plugin }: { plugin: PluginDefinition }) {
  const ecosystem = plugin.id === 'sdk' ? 'TypeScript' : plugin.name;

  return (
    <section className="bg-black py-24 sm:py-32">
      <MarketingContainer>
        <div className="mx-auto max-w-2xl text-center">
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            className="mt-4 text-4xl font-semibold leading-[0.94] tracking-[-0.06em] sm:text-6xl"
          >
            <span className="block">Works with the</span>
            <span className="block text-white/45">{ecosystem} ecosystem.</span>
          </motion.h2>
          <p className="mt-5 text-lg leading-8 text-white/45">
            {plugin.stackDescription}
          </p>
        </div>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5, delay: 0.1 }}
          className="mx-auto mt-14 flex max-w-5xl flex-wrap justify-center gap-3"
        >
          {plugin.technologies.map(({ label, icon: Icon }, index) => (
            <motion.div
              key={label}
              initial={{ opacity: 0, scale: 0.9 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ duration: 0.3, delay: index * 0.04 }}
              className="flex min-h-28 w-[calc(50%_-_0.375rem)] flex-col items-center justify-center gap-3 rounded-2xl border border-white/10 bg-[#080808] px-4 text-center transition-colors hover:border-indigo-300/30 hover:bg-indigo-300/[0.06] sm:w-[calc(33.333%_-_0.5rem)] lg:w-[calc(16.666%_-_0.625rem)]"
            >
              <Icon className="size-8 text-white/70" />
              <span className="text-xs text-white/50">{label}</span>
            </motion.div>
          ))}
        </motion.div>
      </MarketingContainer>
    </section>
  );
}
