import {AbsoluteFill, useCurrentFrame, interpolate} from 'remotion';
import {COLORS, FONTS} from '../styles';

interface BoxProps {
  label: string;
  x: number;
  y: number;
  width: number;
  height: number;
  color: string;
  opacity: number;
}

const Box: React.FC<BoxProps> = ({label, x, y, width, height, color, opacity}) => (
  <div
    style={{
      position: 'absolute',
      left: x,
      top: y,
      width,
      height,
      backgroundColor: COLORS.surface,
      border: `2px solid ${color}`,
      borderRadius: 12,
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      opacity,
    }}
  >
    <span style={{color, fontSize: 28, fontFamily: FONTS.body, fontWeight: 600}}>
      {label}
    </span>
  </div>
);

const Arrow: React.FC<{x1: number; y1: number; x2: number; y2: number; opacity: number}> = ({
  x1, y1, x2, y2, opacity,
}) => (
  <svg
    style={{position: 'absolute', left: 0, top: 0, width: '100%', height: '100%', opacity}}
    viewBox="0 0 1920 1080"
  >
    <defs>
      <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="10" refY="3.5" orient="auto">
        <polygon points="0 0, 10 3.5, 0 7" fill={COLORS.primary} />
      </marker>
    </defs>
    <line
      x1={x1} y1={y1} x2={x2} y2={y2}
      stroke={COLORS.primary}
      strokeWidth={3}
      markerEnd="url(#arrowhead)"
    />
  </svg>
);

export const ArchScene: React.FC = () => {
  const frame = useCurrentFrame();

  const titleOp = interpolate(frame, [0, 30], [0, 1], {extrapolateRight: 'clamp'});
  const box1Op = interpolate(frame, [30, 60], [0, 1], {extrapolateLeft: 'clamp', extrapolateRight: 'clamp'});
  const arrow1Op = interpolate(frame, [60, 90], [0, 1], {extrapolateLeft: 'clamp', extrapolateRight: 'clamp'});
  const box2Op = interpolate(frame, [90, 120], [0, 1], {extrapolateLeft: 'clamp', extrapolateRight: 'clamp'});
  const arrow2Op = interpolate(frame, [120, 150], [0, 1], {extrapolateLeft: 'clamp', extrapolateRight: 'clamp'});
  const box3Op = interpolate(frame, [150, 180], [0, 1], {extrapolateLeft: 'clamp', extrapolateRight: 'clamp'});
  const arrow3Op = interpolate(frame, [180, 210], [0, 1], {extrapolateLeft: 'clamp', extrapolateRight: 'clamp'});
  const box4Op = interpolate(frame, [210, 240], [0, 1], {extrapolateLeft: 'clamp', extrapolateRight: 'clamp'});

  return (
    <AbsoluteFill style={{backgroundColor: COLORS.bg, padding: 80}}>
      <h1
        style={{
          color: COLORS.text,
          fontSize: 64,
          fontFamily: FONTS.heading,
          fontWeight: 700,
          textAlign: 'center',
          opacity: titleOp,
          marginBottom: 40,
        }}
      >
        How it works
      </h1>

      <div style={{position: 'relative', flex: 1}}>
        <Box label="Capture" x={160} y={200} width={240} height={100} color={COLORS.primary} opacity={box1Op} />
        <Arrow x1={400} y1={250} x2={560} y2={250} opacity={arrow1Op} />
        <Box label="Triage" x={560} y={200} width={240} height={100} color={COLORS.accent} opacity={box2Op} />
        <Arrow x1={800} y1={250} x2={960} y2={250} opacity={arrow2Op} />
        <Box label="Build + Review" x={960} y={200} width={300} height={100} color={COLORS.warning} opacity={box3Op} />
        <Arrow x1={1260} y1={250} x2={1420} y2={250} opacity={arrow3Op} />
        <Box label="Archive + Present" x={1420} y={200} width={320} height={100} color={COLORS.primary} opacity={box4Op} />

        {/* Data flow labels */}
        <div
          style={{
            position: 'absolute',
            left: 160,
            top: 380,
            opacity: interpolate(frame, [240, 270], [0, 1], {
              extrapolateLeft: 'clamp',
              extrapolateRight: 'clamp',
            }),
          }}
        >
          {[
            '1. User input → REQ files + UR folder (verbatim preserved)',
            '2. Triage by complexity → Route A / B / C',
            '3. Plan → Explore → Implement → Test → Review → Commit',
            '4. Archive with full history → Client-ready deliverables',
          ].map((step, i) => (
            <div
              key={i}
              style={{
                color: COLORS.muted,
                fontSize: 30,
                fontFamily: FONTS.body,
                marginBottom: 16,
                opacity: interpolate(frame, [260 + i * 30, 280 + i * 30], [0, 1], {
                  extrapolateLeft: 'clamp',
                  extrapolateRight: 'clamp',
                }),
              }}
            >
              {step}
            </div>
          ))}
        </div>
      </div>
    </AbsoluteFill>
  );
};
