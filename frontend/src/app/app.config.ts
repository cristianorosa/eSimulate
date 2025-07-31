import { ApplicationConfig, importProvidersFrom } from "@angular/core";
import { provideRouter } from "@angular/router";
import {
  provideHttpClient,
  withInterceptors,
  withXsrfConfiguration,
} from "@angular/common/http";
import { provideAnimations } from "@angular/platform-browser/animations";
import { provideZoneChangeDetection } from "@angular/core";
import { routes } from "./app.routes";
import { AuthInterceptor } from "./core/auth.interceptor";

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideAnimations(),
    provideZoneChangeDetection(),
    provideHttpClient(
      withInterceptors([AuthInterceptor]),
      withXsrfConfiguration({
        cookieName: "XSRF-TOKEN",
        headerName: "X-XSRF-TOKEN",
      })
    ),
    // Outros providers globais...
  ],
};
