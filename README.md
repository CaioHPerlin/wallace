# WallAce

### Produto: Sistema de Alarmes Remotamente Configurável

### Integrantes: Caio Hygino Perlin, Caio Ricarte Loyer e Murilo Toledo

## Objetivo do Projeto

Criar um sistema IoT de detecção e monitoramento de movimento utilizando um ESP32 e um sensor PIR, permitindo o acompanhamento das detecções em tempo real por meio de um dashboard web e o controle remoto dos feedbacks físicos do dispositivo.

## A Necessidade (Por que construir?)

Ambientes residenciais, comerciais ou institucionais podem necessitar de uma forma simples de identificar movimentações e acompanhar esses eventos remotamente. Soluções tradicionais de monitoramento podem ser mais complexas, caras ou pouco flexíveis para projetos de pequeno porte. O WallAce busca oferecer uma alternativa simples e conectada, capaz de registrar detecções de movimento, disponibilizar essas informações em tempo real e permitir que o usuário controle e configure remotamente os alertas físicos do dispositivo.

## Descrição do Produto (Como vai funcionar)

O WallAce utilizará um dev kit ESP32 conectado à rede Wi-Fi. Na camada de percepção, será utilizado um sensor PIR, responsável por detectar movimentações no ambiente. O dispositivo também contará com um LED e um buzzer como atuadores para fornecer feedback visual e sonoro sempre que uma movimentação for identificada.

- Edge/Borda: O ESP32 realizará continuamente a leitura do sensor PIR. Ao detectar movimento, o dispositivo registrará o evento, poderá acionar localmente o LED e o buzzer e enviará a detecção ao servidor através da rede Wi-Fi. O funcionamento dos atuadores poderá ser habilitado ou desabilitado remotamente pelo usuário, sem necessidade de reprogramar o ESP32.
- Nuvem/Aplicação: O backend será desenvolvido em Go e ficará responsável por receber os eventos enviados pelo ESP32 e distribuí-los em tempo real para o dashboard. A comunicação entre o servidor e a interface web será realizada por WebSocket. O frontend, desenvolvido em TypeScript, exibirá as detecções recentes e manterá um histórico dos eventos registrados. Pelo dashboard, o usuário também poderá habilitar ou desabilitar remotamente o buzzer e o LED, e o servidor enviará essas configurações de volta ao ESP32.

<img width="669" height="518" alt="image" src="https://github.com/user-attachments/assets/e7d5ecc3-24f2-4bd4-b088-9cb4fad7cf0e" />
