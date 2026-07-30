import { CustomDomainSection } from './custom-domain';
import { DeveloperExperience } from './developer-experience';
import { Hero } from './hero';
import { LocalAccessSection } from './local-access';
import { NetworkDiagram } from './network-diagram';
import { OpenSourceSection } from './open-source';
import { ProtocolsSection } from './protocols';

export function LandingPage() {
  return (
    <>
      <Hero />
      <DeveloperExperience />
      <NetworkDiagram />
      <CustomDomainSection />
      <LocalAccessSection />
      <ProtocolsSection />
      <OpenSourceSection />
    </>
  );
}
