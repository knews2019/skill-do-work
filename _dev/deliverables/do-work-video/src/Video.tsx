import {AbsoluteFill, Sequence} from 'remotion';
import {ProblemScene} from './scenes/ProblemScene';
import {SolutionScene} from './scenes/SolutionScene';
import {ArchScene} from './scenes/ArchScene';
import {ValueScene} from './scenes/ValueScene';
import {COLORS} from './styles';

export const Video: React.FC = () => {
  return (
    <AbsoluteFill style={{backgroundColor: COLORS.bg}}>
      <Sequence from={0} durationInFrames={450}>
        <ProblemScene />
      </Sequence>
      <Sequence from={450} durationInFrames={900}>
        <SolutionScene />
      </Sequence>
      <Sequence from={1350} durationInFrames={600}>
        <ArchScene />
      </Sequence>
      <Sequence from={1950} durationInFrames={450}>
        <ValueScene />
      </Sequence>
    </AbsoluteFill>
  );
};
