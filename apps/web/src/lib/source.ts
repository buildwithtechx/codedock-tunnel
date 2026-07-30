import { docs } from 'fumadocs-mdx:collections/server';
import { loader } from 'fumadocs-core/source';
import { icons } from 'lucide-react';
import { type ComponentType, createElement } from 'react';
import {
  SiExpress,
  SiNestjs,
  SiNextdotjs,
  SiReact,
  SiTypescript,
  SiVite,
} from 'react-icons/si';

const customIcons: Record<string, ComponentType<{ className?: string }>> = {
  SiExpress,
  SiNestjs,
  SiNextdotjs,
  SiReact,
  SiTypescript,
  SiVite,
};

export const source = loader({
  baseUrl: '/docs',
  source: docs.toFumadocsSource(),
  icon(icon) {
    if (!icon) return;
    if (icon in customIcons) return createElement(customIcons[icon]);
    if (icon in icons) return createElement(icons[icon as keyof typeof icons]);
  },
});
