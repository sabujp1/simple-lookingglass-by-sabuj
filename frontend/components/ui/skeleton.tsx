import * as React from 'react';
import { cn } from '@/lib/utils';

function getRandomDelay() {
  return Math.floor(Math.random() * 400) + 100;
}

interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {
  delay?: number;
}

const Skeleton = React.forwardRef<HTMLDivElement, SkeletonProps>(
  ({ className, delay, style, ...props }, ref) => {
    const [visible, setVisible] = React.useState(false);
    
    React.useEffect(() => {
      const timer = setTimeout(() => setVisible(true), delay ?? 0);
      return () => clearTimeout(timer);
    }, [delay]);

    return (
      <div
        ref={ref}
        className={cn(
          'animate-pulse rounded-md bg-muted',
          visible ? 'opacity-100' : 'opacity-0',
          className
        )}
        style={{
          transition: 'opacity 300ms ease-in-out',
          ...style,
        }}
        {...props}
      />
    );
  }
);
Skeleton.displayName = 'Skeleton';

export { Skeleton };