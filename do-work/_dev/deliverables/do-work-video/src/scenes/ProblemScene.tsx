import {AbsoluteFill, useCurrentFrame, interpolate, Sequence} from 'remotion';
import {COLORS, FONTS} from '../styles';

const PAIN_POINTS = [
  'Every conversation starts from scratch',
  'Complex requests get half-built',
  'No paper trail — no way to trace what happened',
  'Client asks "what did I get?" — you scroll chat logs',
];

export const ProblemScene: React.FC = () => {
  const frame = useCurrentFrame();

  const titleOpacity = interpolate(frame, [0, 30], [0, 1], {
    extrapolateRight: 'clamp',
  });
  const titleY = interpolate(frame, [0, 30], [20, 0], {
    extrapolateRight: 'clamp',
  });

  return (
    <AbsoluteFill
      style={{
        backgroundColor: COLORS.surface,
        justifyContent: 'center',
        alignItems: 'center',
        padding: 100,
      }}
    >
      <div
        style={{
          opacity: titleOpacity,
          transform: `translateY(${titleY}px)`,
          marginBottom: 60,
        }}
      >
        <h1
          style={{
            color: COLORS.text,
            fontSize: 80,
            fontFamily: FONTS.heading,
            fontWeight: 700,
            textAlign: 'center',
            margin: 0,
          }}
        >
          AI coding tools have a memory problem
        </h1>
      </div>

      <div style={{display: 'flex', flexDirection: 'column', gap: 24}}>
        {PAIN_POINTS.map((point, i) => {
          const startFrame = 60 + i * 30;
          const opacity = interpolate(frame, [startFrame, startFrame + 20], [0, 1], {
            extrapolateLeft: 'clamp',
            extrapolateRight: 'clamp',
          });
          const x = interpolate(frame, [startFrame, startFrame + 20], [-40, 0], {
            extrapolateLeft: 'clamp',
            extrapolateRight: 'clamp',
          });

          return (
            <div
              key={i}
              style={{
                opacity,
                transform: `translateX(${x}px)`,
                display: 'flex',
                alignItems: 'center',
                gap: 16,
              }}
            >
              <div
                style={{
                  width: 12,
                  height: 12,
                  borderRadius: 6,
                  backgroundColor: COLORS.warning,
                  flexShrink: 0,
                }}
              />
              <span
                style={{
                  color: COLORS.muted,
                  fontSize: 40,
                  fontFamily: FONTS.body,
                }}
              >
                {point}
              </span>
            </div>
          );
        })}
      </div>
    </AbsoluteFill>
  );
};
