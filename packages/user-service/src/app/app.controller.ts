import { Controller, Request, Get } from '@nestjs/common';

import { AppService } from './app.service';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get()
  getData(@Request() req: Request) {
    return this.appService.getData(req);
  }
}
