import { Github, ExternalLink, Heart, Code, FileJson } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';

export default function About() {
  return (
    <div className="p-6 max-w-3xl mx-auto">
      <div className="text-center mb-8">
        <div className="flex justify-center mb-4">
          <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-primary text-primary-foreground font-mono font-bold text-2xl">
            MM
          </div>
        </div>
        <h1 className="text-3xl font-bold">MMock Console</h1>
        <p className="text-muted-foreground mt-2">
          A modern web interface for MMock HTTP mock server
        </p>
        <Badge variant="secondary" className="mt-3 font-mono">
          v1.0.0
        </Badge>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Code className="h-5 w-5" />
              About MMock
            </CardTitle>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground space-y-3">
            <p>
              MMock is a lightweight, configurable HTTP mock server designed for testing
              and rapid API prototyping. It allows developers to simulate API endpoints
              during development without requiring a real backend.
            </p>
            <p>
              Key features include request matching with wildcards, dynamic response
              generation, state machine scenarios, and real-time request monitoring.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FileJson className="h-5 w-5" />
              Console Features
            </CardTitle>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground">
            <ul className="space-y-2">
              <li className="flex items-center gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-primary flex-shrink-0" />
                Real-time request monitoring via WebSocket
              </li>
              <li className="flex items-center gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-primary flex-shrink-0" />
                Create, edit, and delete mock definitions
              </li>
              <li className="flex items-center gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-primary flex-shrink-0" />
                Filter requests by match status
              </li>
              <li className="flex items-center gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-primary flex-shrink-0" />
                JSON syntax highlighting
              </li>
              <li className="flex items-center gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-primary flex-shrink-0" />
                Scenario state management
              </li>
            </ul>
          </CardContent>
        </Card>
      </div>

      <Separator className="my-8" />

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Heart className="h-5 w-5 text-destructive" />
            Credits
          </CardTitle>
          <CardDescription>
            MMock is an open-source project created and maintained by the community
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">Author</p>
              <p className="text-sm text-muted-foreground">Jordi Martin</p>
            </div>
            <Button variant="outline" asChild>
              <a href="https://github.com/jmartin82" target="_blank" rel="noopener noreferrer">
                <Github className="h-4 w-4 mr-2" />
                GitHub Profile
              </a>
            </Button>
          </div>

          <Separator />

          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">License</p>
              <p className="text-sm text-muted-foreground">MIT License</p>
            </div>
            <Badge variant="outline">Open Source</Badge>
          </div>

          <Separator />

          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div>
              <p className="font-medium">MMock Repository</p>
              <p className="text-sm text-muted-foreground">View source code and documentation</p>
            </div>
            <Button asChild>
              <a href="https://github.com/jmartin82/mmock" target="_blank" rel="noopener noreferrer">
                <Github className="h-4 w-4 mr-2" />
                View on GitHub
                <ExternalLink className="h-3 w-3 ml-2" />
              </a>
            </Button>
          </div>
        </CardContent>
      </Card>

      <div className="text-center mt-8 text-sm text-muted-foreground">
        <p>Made with ❤️ for developers who need reliable API mocking</p>
      </div>
    </div>
  );
}
