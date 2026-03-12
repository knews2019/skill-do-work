import {Composition} from 'remotion';
import {Video} from './Video';

export const RemotionRoot: React.FC = () => {
  return (
    <Composition
      id="DoWork"
      component={Video}
      durationInFrames={2700}
      fps={30}
      width={1920}
      height={1080}
    />
  );
};
