import * as AvatarPrimitive from '@radix-ui/react-avatar'
import type * as React from 'react'
import { cn } from '@/lib/utils'

type AvatarRootProps = React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Root>
type AvatarImageProps = React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Image>
type AvatarFallbackProps = React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Fallback>

export function Avatar({ className, ...props }: AvatarRootProps) {
  return <AvatarPrimitive.Root className={cn('relative flex h-8 w-8 shrink-0 overflow-hidden rounded-full', className)} {...props} />
}

export function AvatarImage({ className, ...props }: AvatarImageProps) {
  return <AvatarPrimitive.Image className={cn('aspect-square h-full w-full', className)} {...props} />
}

export function AvatarFallback({ className, ...props }: AvatarFallbackProps) {
  return <AvatarPrimitive.Fallback className={cn('flex h-full w-full items-center justify-center rounded-full bg-muted text-xs', className)} {...props} />
}
