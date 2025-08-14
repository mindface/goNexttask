# OpenAPI Specification — Go Nexttask (YAML)

```yaml
openapi: 3.0.3
info:
  title: Go Nexttask API
  version: 1.0.0
  description: ベアリング製造・金属加工向け生産管理、品質管理、NC加工連携システムAPI

servers:
  - url: https://api.gonexttask.local/v1

paths:
  /production-plans:
    get:
      summary: 生産計画一覧取得
      responses:
        '200':
          description: 生産計画一覧
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/ProductionPlan'

    post:
      summary: 生産計画作成
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ProductionPlanCreate'
      responses:
        '201':
          description: 作成成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ProductionPlan'

  /production-plans/{id}:
    get:
      summary: 生産計画詳細取得
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: 生産計画詳細
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ProductionPlan'

  /nc-programs:
    get:
      summary: NCプログラム一覧取得
      responses:
        '200':
          description: NCプログラム一覧
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/NcProgram'

    post:
      summary: NCプログラム登録
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/NcProgramCreate'
      responses:
        '201':
          description: 登録成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/NcProgram'

  /quality/inspections:
    post:
      summary: 検査結果登録
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/InspectionCreate'
      responses:
        '201':
          description: 登録成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Inspection'

components:
  schemas:
    ProductionPlan:
      type: object
      properties:
        id:
          type: string
          example: "plan-1234"
        orderId:
          type: string
          example: "order-5678"
        material:
          type: string
          example: "Steel A"
        quantity:
          type: integer
          example: 100
        status:
          type: string
          enum: [planned, in_progress, completed, delayed]
          example: planned
        scheduledStartDate:
          type: string
          format: date-time
        scheduledEndDate:
          type: string
          format: date-time

    ProductionPlanCreate:
      type: object
      required:
        - orderId
        - material
        - quantity
        - scheduledStartDate
        - scheduledEndDate
      properties:
        orderId:
          type: string
        material:
          type: string
        quantity:
          type: integer
        scheduledStartDate:
          type: string
          format: date-time
        scheduledEndDate:
          type: string
          format: date-time

    NcProgram:
      type: object
      properties:
        id:
          type: string
          example: "ncprog-001"
        partId:
          type: string
          example: "part-100"
        machineId:
          type: string
          example: "machine-01"
        version:
          type: string
          example: "v1.0.3"
        data:
          type: string
          description: "NCプログラムの内容（ベース64エンコードなど）"

    NcProgramCreate:
      type: object
      required:
        - partId
        - machineId
        - version
        - data
      properties:
        partId:
          type: string
        machineId:
          type: string
        version:
          type: string
        data:
          type: string

    Inspection:
      type: object
      properties:
        id:
          type: string
          example: "inspection-777"
        lotNumber:
          type: string
          example: "lot-20250811-001"
        machineId:
          type: string
          example: "machine-01"
        operatorId:
          type: string
          example: "operator-123"
        result:
          type: string
          enum: [pass, fail]
        measuredValues:
          type: object
          additionalProperties:
            type: number
        inspectionDate:
          type: string
          format: date-time

    InspectionCreate:
      type: object
      required:
        - lotNumber
        - machineId
        - operatorId
        - result
        - measuredValues
        - inspectionDate
      properties:
        lotNumber:
          type: string
        machineId:
          type: string
        operatorId:
          type: string
        result:
          type: string
          enum: [pass, fail]
        measuredValues:
          type: object
          additionalProperties:
            type: number
        inspectionDate:
          type: string
          format: date-time
```
