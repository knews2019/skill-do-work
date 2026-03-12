import {AbsoluteFill, useCurrentFrame, interpolate, spring, useVideoConfig} from 'remotion';
import {COLORS, FONTS} from '../styles';

const VALUES = [
  {before: 'Chat logs lost between sessions', after: 'Persistent archive with full history'},
  {before: 'Manual status updates for clients', after: 'One-command client briefs + video'},
  {before: '"wip" commits, no paper trail', after: 'Atomic commits per request with review scores'},
  {before: 'Context lost when switching AI tools', after: 'Portable archive — tool-agnostic'},
];

export const ValueScene: React.FC = () => {
  const frame = useCurrentFrame();
  const {fps} = useVideoConfig();

  const titleOp = interpolate(frame, [0, 30], [0, 1], {extrapolateRight: 'clamp'});

  const ctaOp = interpolate(frame, [330, 360], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
  });
  const ctaScale = spring({
    frame: frame - 330,
    fps,
    config: {damping: 10, stiffness: 100},
  });

  return (
    <AbsoluteFill
      style={{
        backgroundColor: COLORS.surface,
        padding: 80,
        justifyContent: 'center',
        alignItems: 'center',
      }}
    >
      <h1
        style={{
          color: COLORS.text,
          fontSize: 64,
          fontFamily: FONTS.heading,
          fontWeight: 700,
          textAlign: 'center',
          opacity: titleOp,
          marginBottom: 60,
        }}
      >
        Before &amp; After
      </h1>

      <div style={{display: 'flex', gap: 60, marginBottom: 60}}>
        {/* Before column */}
        <div style={{flex: 1}}>
          <h2 style={{color: COLORS.warning, fontSize: 36, fontFamily: FONTS.heading, marginBottom: 24}}>
            Without do-work
          </h2>
          {VALUES.map((v, i) => {
            const op = interpolate(frame, [40 + i * 40, 60 + i * 40], [0, 1], {
              extrapolateLeft: 'clamp',
              extrapolateRight: 'clamp',
            });
            return (
              <div key={i} style={{opacity: op, marginBottom: 20}}>
                <span style={{color: COLORS.muted, fontSize: 28, fontFamily: FONTS.body}}>
                  {v.before}
                </span>
              </div>
            );
          })}
        </div>

        {/* After column */}
        <div style={{flex: 1}}>
          <h2 style={{color: COLORS.accent, fontSize: 36, fontFamily: FONTS.heading, marginBottom: 24}}>
            With do-work
          </h2>
          {VALUES.map((v, i) => {
            const op = interpolate(frame, [60 + i * 40, 80 + i * 40], [0, 1], {
              extrapolateLeft: 'clamp',
              extrapolateRight: 'clamp',
            });
            return (
              <div key={i} style={{opacity: op, marginBottom: 20}}>
                <span style={{color: COLORS.text, fontSize: 28, fontFamily: FONTS.body, fontWeight: 600}}>
                  {v.after}
                </span>
              </div>
            );
          })}
        </div>
      </div>

      {/* CTA */}
      <div
        style={{
          opacity: ctaOp,
          transform: `scale(${ctaScale})`,
          textAlign: 'center',
        }}
      >
        <code
          style={{
            color: COLORS.accent,
            fontSize: 40,
            fontFamily: FONTS.mono,
            backgroundColor: COLORS.bg,
            padding: '16px 40px',
            borderRadius: 12,
            border: `1px solid ${COLORS.border}`,
          }}
        >
          npx skills add knews2019/skill-do-work
        </code>
        <p style={{color: COLORS.muted, fontSize: 28, fontFamily: FONTS.body, marginTop: 20}}>
          One install. Seven actions. Every task captured, built, reviewed, and ready to present.
        </p>
      </div>
    </AbsoluteFill>
  );
};
