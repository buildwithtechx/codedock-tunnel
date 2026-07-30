import { motion } from 'motion/react';
import { MarketingContainer } from '#/components/layout';
import type { PluginDefinition } from './plugin-data';

export function PluginStackSection({ plugin }: { plugin: PluginDefinition }) {
  return (
    <section className="bg-white/[0.012] py-24 sm:py-32">
      <MarketingContainer>
        <div className="mx-auto max-w-2xl text-center">
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            className="mt-4 text-4xl font-semibold tracking-[-0.06em] sm:text-6xl"
          >
            {plugin.stackHeading}
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
          className="mx-auto mt-14 grid max-w-5xl grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6"
        >
          {plugin.technologies.map(({ label, icon: Icon }) => (
            <motion.div
              key={label}
              initial={{ opacity: 0, scale: 0.9 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ duration: 0.3 }}
              className="flex min-h-28 flex-col items-center justify-center gap-3 rounded-2xl border border-white/10 bg-[#0a0a0a] px-4 text-center transition-colors hover:border-indigo-300/30 hover:bg-indigo-300/[0.06]"
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
