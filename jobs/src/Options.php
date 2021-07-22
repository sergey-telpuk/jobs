<?php

/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

declare(strict_types=1);

namespace Spiral\Jobs;

final class Options implements \JsonSerializable
{
    /** @var int|null */
    private $delay = null;

    /** @var string|null */
    private $pipeline = null;

    /**
     * @param int $delay
     * @return self
     */
    public function withDelay(?int $delay): self
    {
        $options = clone $this;
        $options->delay = $delay;

        return $options;
    }

    /**
     * @return string|null
     */
    public function getPipeline(): ?string
    {
        return $this->pipeline;
    }

    /**
     * @param string|null $pipeline
     * @return self
     */
    public function withPipeline(?string $pipeline): self
    {
        $options = clone $this;
        $options->pipeline = $pipeline;

        return $options;
    }

    /**
     * @return int|null
     */
    public function getDelay(): ?int
    {
        return $this->delay;
    }

    /**
     * @return array|mixed
     */
    public function jsonSerialize()
    {
        return [
            'delay'    => $this->delay,
            'pipeline' => $this->pipeline
        ];
    }

    /**
     * @param int $delay
     * @return Options
     */
    public static function delayed(int $delay): Options
    {
        $options = new self();
        $options->delay = $delay;

        return $options;
    }
}
