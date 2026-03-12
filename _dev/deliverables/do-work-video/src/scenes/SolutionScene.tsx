import {AbsoluteFill, useCurrentFrame, interpolate, Sequence, spring, useVideoConfig} from 'remotion';
import {COLORS, FONTS} from '../styles';

const STEPS = [
  {cmd: 'do work add dark mode, fix search, align header', label: 'Capture — 3 tasks from 1 sentence'},
  {cmd: 'do work run', label: 'Build — triage, plan, implement, test, review, commit'},
  {cmd: 'do work present work', label: 'Present — client brief, architecture, Remotion video'},
];

export const SolutionScene: React.FC = () => {
  const frame = useCurrentFrame();
  const {fps} = useVideoConfig();

  const titleOpacity = interpolate(frame, [0, 30], [0, 1], {
    extrapolateRight: 'clamp',
  });

  return (
    <AbsoluteFill
      style={{
        backgroundColor: COLORS.bg,
        padding: 100,
        justifyContent: 'center',
      }}
    >
      <h1
        style={{
          color: COLORS.primary,
          fontSize: 64,
          fontFamily: FONTS.heading,
          fontWeight: 700,
          opacity: titleOpacity,
          marginBottom: 80,
          textAlign: 'center',
        }}
      >
        Separate thinking from doing
      </h1>

      {STEPS.map((step, i) => {
        const enterFrame = 60 + i * 200;
        const scale = spring({
          frame: frame - enterFrame,
          fps,
          config: {damping: 12, stiffness: 80},
        });
        const opacity = interpolate(frame, [enterFrame, enterFrame + 20], [0, 1], {
          extrapolateLeft: 'clamp',
          extrapolateRight: 'clamp',
        });

        return (
          <div
            key={i}
            style={{
              opacity,
              transform: `scale(${scale})`,
              marginBottom: 48,
            }}
          >
            <div
              style={{
                backgroundColor: COLORS.surface,
                borderRadius: 16,
                padding: '32px 48px',
                border: `1px solid ${COLORS.border}`,
              }}
            >
              <code
                style={{
                  color: COLORS.accent,
                  fontSize: 36,
                  fontFamily: FONTS.mono,
                  display: 'block',
                  marginBottom: 12,
                }}
              >
                $ {step.cmd}
              </code>
              <span
                style={{
                  color: COLORS.muted,
                  fontSize: 28,
                  fontFamily: FONTS.body,
                }}
              >
                {step.label}
              </span>
            </div>
          </div>
        );
      })}
    </AbsoluteFill>
  );
};
